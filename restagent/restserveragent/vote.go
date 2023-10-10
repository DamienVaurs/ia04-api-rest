package restserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

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
		if list[i] != comsoc.Alternative(i+1) {
			return false
		}
	}
	return true
}

func checkVote(ballotsList map[string]restagent.Ballot, req restagent.RequestVote) (err error) {
	//Vérifie que le ballot existe
	_, found := ballotsList[req.BallotId]
	if !found {
		return fmt.Errorf("notexist")
	}
	//vérifie que l'agent n'a pas déjà voté
	fmt.Println("Ont voté  : ", ballotsList[req.BallotId].HaveVoted)
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

	//Vérifie que la date de fin est n'est pas passée
	/*if rsa.ballotsList[req.BallotId].Deadline.Before(time.Now()) {
		return fmt.Errorf("alreadyfinished")
	}*/
	//Vérifie que les alteratives fournies pour le vote sont correctes
	if !checkVoteAlts(req.Prefs, ballotsList[req.BallotId].Alts) {
		return fmt.Errorf("wrongalts")
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Println("Serveur recoit : ", r.URL, req)

	//Vérifie que le vote est correct
	err = checkVote(rsa.ballotsList, req)
	if err != nil {
		switch err.Error() {
		case "notexist":
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /vote : ballot %s does not exist", req.BallotId)
			w.Write([]byte(msg))
			return
		case "alreadyvoted":
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /vote : agent %s has already voted for ballot %s", req.AgentId, req.BallotId)
			w.Write([]byte(msg))
			return
		case "notallowed":
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /vote : agent %s is not allowed to vote for ballot %s", req.AgentId, req.BallotId)
			w.Write([]byte(msg))
			return
		case "alreadyfinished":
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /vote : ballot %s is already finished : %s", req.BallotId, rsa.ballotsList[req.BallotId].Deadline.String())
			w.Write([]byte(msg))
			return
		case "wrongalts":
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /vote : alternatives provided for ballot %s are not correct", req.BallotId)
			w.Write([]byte(msg))
			return
		default:
			fmt.Println("Vote correct ", req)
		}
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

	w.WriteHeader(http.StatusOK)
	serial, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Write(serial)
}
