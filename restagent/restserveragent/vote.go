package restserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
)

// Fonctions qui traitent l'appel à l'API REST pour voter:
// http://localhost:8080/vote

// Décode la requête
func (*RestServerAgent) decodeVoteRequest(r *http.Request) (req restagent.RequestVote, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		fmt.Println("Erreur de décodage de la requête /vote : ", err)

	}
	return
}

func checkVoteAlts(vote []comsoc.Alternative, expected int) bool {
	//vérifie que le vote correspond aux alternatives proposées par le ballot

	//Rappel : les alternatives attendues sont entre 1 et Ballot.Alts inclus
	if len(vote) != expected {
		return false
	}
	list := make([]comsoc.Alternative, expected)
	copy(list, vote)
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	for i := 0; i < expected-1; i++ {
		if list[i] != list[i+1]-1 {
			return false
		}
	}
	return true
}

func checkVote(ballotsList map[string]restagent.Ballot, deadline time.Time, req restagent.RequestVote) (err error) {
	//Vérifie que le ballot existe
	_, found := ballotsList[req.BallotId]
	if !found {
		return fmt.Errorf("notexist")
	}
	//vérifie que l'agent n'a pas déjà voté
	for _, v := range ballotsList[req.BallotId].HaveVoted {
		if v == req.AgentId {
			return fmt.Errorf("alreadyvoted")
		}
	}
	//Vérifie que l'agent a le droit de voter
	var canVote bool
	for _, v := range ballotsList[req.BallotId].VoterIds {
		if v == req.AgentId {
			canVote = true
			break
		}
	}
	if !canVote {
		return fmt.Errorf("notallowed")
	}

	//Vérifie que la date de clôture n'est pas passée
	if deadline.Before(time.Now()) {
		return fmt.Errorf("alreadyfinished")
	}

	//Vérifie que les alteratives fournies pour le vote sont correctes
	if !checkVoteAlts(req.Prefs, ballotsList[req.BallotId].Alts) {
		return fmt.Errorf("wrongalts")
	}

	//Si le ballot est "approval" vérifie qu'un seuil cohérent est bien fourni
	//Remarque : discuter dans le README de ce choix, car on aurait pu imaginer que pas de treshold => on compte tout le monde
	if ballotsList[req.BallotId].Rule == "approval" {
		if req.Options == nil || len(req.Options) != 1 || req.Options[0] < 0 || req.Options[0] > ballotsList[req.BallotId].Alts {
			//TODO : Vérifier si Threshold commence à 0 ou à 1 pour valider ce test
			return fmt.Errorf("wrongthreshold")
		}
	}
	return nil
}

func (rsa *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()
	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeVoteRequest(r)
	if err != nil {
		//TODO : bérifier si c'est bien le bon code d'erreur
		w.WriteHeader(http.StatusInternalServerError) //500
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Println("Serveur recoit : ", r.URL, req)

	//Vérifie que le vote est correct
	err = checkVote(rsa.ballotsList, rsa.ballotsList[req.BallotId].Deadline, req)
	if err != nil {
		switch err.Error() {
		case "notexist":
			w.WriteHeader(http.StatusBadRequest) //400
			msg := fmt.Sprintf("error /vote : ballot %s does not exist", req.BallotId)
			w.Write([]byte(msg))
			return
		case "alreadyvoted":
			w.WriteHeader(http.StatusForbidden) //403
			msg := fmt.Sprintf("error /vote : agent %s has already voted for ballot %s", req.AgentId, req.BallotId)
			w.Write([]byte(msg))
			return
		case "notallowed":
			w.WriteHeader(http.StatusUnauthorized) //401
			msg := fmt.Sprintf("error /vote : agent %s is not allowed to vote for ballot %s", req.AgentId, req.BallotId)
			w.Write([]byte(msg))
			return
		case "alreadyfinished":
			w.WriteHeader(http.StatusServiceUnavailable) //503
			msg := fmt.Sprintf("error /vote : ballot %s is already finished : %s", req.BallotId, rsa.ballotsList[req.BallotId].Deadline.String())
			w.Write([]byte(msg))
			return
		case "wrongalts":
			w.WriteHeader(http.StatusBadRequest) //400
			msg := fmt.Sprintf("error /vote : alternatives provided for ballot %s are not correct", req.BallotId)
			w.Write([]byte(msg))
			return
		case "wrongthreshold":
			w.WriteHeader(http.StatusBadRequest) //400
			msg := fmt.Sprintf("error /vote : threshold %d provided for ballot %s is not correct", req.Options, req.BallotId)
			w.Write([]byte(msg))
			return

		}
	}

	//Enregistre le threshold si besoin
	if rsa.ballotsList[req.BallotId].Rule == restagent.Approval {
		_, found := rsa.ballotsList[req.BallotId].Thresholds[req.AgentId]
		if found {
			w.WriteHeader(http.StatusBadRequest) //400
			msg := fmt.Sprintf("error /vote : agent %s has already provided a threshold for ballot %s", req.AgentId, req.BallotId)
			w.Write([]byte(msg))
			return
		}
		rsa.ballotsList[req.BallotId].Thresholds[req.AgentId] = req.Options[0]
	}

	//Enregistre le vote pour le ballot
	rsa.ballotsMap[req.BallotId] = append(rsa.ballotsMap[req.BallotId], req.Prefs)

	//Enregistre que l'agent a voté
	for i := 0; i < len(rsa.ballotsList[req.BallotId].HaveVoted); i++ {
		if rsa.ballotsList[req.BallotId].HaveVoted[i] == "" {
			rsa.ballotsList[req.BallotId].HaveVoted[i] = req.AgentId
			break
		}
	}

	w.WriteHeader(http.StatusOK) //200
	msg := "/vote : vote registered"
	w.Write([]byte(msg))
}
