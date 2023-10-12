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

// Fonctions qui traitent l'appel à l'API REST pour créer un ballot:
// http://localhost:8080/new_ballot

// Décode la requête
func (*RestServerAgent) decodeNewBallotRequest(r *http.Request) (req restagent.RequestNewBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		fmt.Println("Erreur de décodage de la requête /new_ballot : ", err)
		return
	}
	return
}

// Effectue plusieurs vérifications sur le scrutin fournit
func checkBallot(req restagent.RequestNewBallot) (err error) {
	//Vérifie que le format de date est bon
	_, err = time.Parse(time.RFC3339, req.Deadline) //A vérifier
	if err != nil {
		return fmt.Errorf("deadline")
	}

	//Vérifie que le type de ballot est autorisé
	var authorized = false
	for _, v := range restagent.Rules {
		if v == req.Rule {
			authorized = true
			break
		}
	}
	if !authorized {
		return fmt.Errorf("rule")
	}

	//Vérifie que les alternatives sont cohérentes avec le tie-break
	if req.Alts < 1 {
		return fmt.Errorf("alts")
	}

	//Vérifie que le tie-break est cohérent avec les alternatives
	if req.TieBreak == nil || len(req.TieBreak) != req.Alts {
		return fmt.Errorf("tiebreak")
	} else {
		//Vérifie qu'il n'y a pas de doublon dans le tie-break ni de valeur abérante
		list := make([]comsoc.Alternative, len(req.TieBreak))
		copy(list, req.TieBreak)
		sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
		for i := 0; i < len(req.TieBreak)-1; i++ {
			if list[i]+1 != list[i+1] {
				return fmt.Errorf("tiebreak")
			}
		}
	}

	return
}

func (rsa *RestServerAgent) doCreateNewBallot(w http.ResponseWriter, r *http.Request) {
	rsa.Lock() // verrouillage de la méthode nécessaire au moins car la méthode vote modifie countBallot. Sinon, on pourrait utiliser une goroutine par ballot
	defer rsa.Unlock()
	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeNewBallotRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Println("Serveur recoit : ", r.URL, req)

	err = checkBallot(req)

	if err != nil {
		switch err.Error() {
		case "deadline":
			w.WriteHeader(http.StatusBadRequest) //400
			msg := fmt.Sprintf("error /new_ballot : deadline %s is not in the right format", req.Deadline)
			w.Write([]byte(msg))
			return
		case "rule":
			w.WriteHeader(http.StatusNotImplemented) //501
			msg := fmt.Sprintf("error /new_ballot : rule %s is not implemented", req.Rule)
			w.Write([]byte(msg))
			return
		case "alts":
			w.WriteHeader(http.StatusBadRequest) //400
			msg := fmt.Sprintf("error /new_ballot : number of alternatives %d should be >=1", req.Alts)
			w.Write([]byte(msg))
			return
		case "tiebreak":
			w.WriteHeader(http.StatusBadRequest) //400
			msg := fmt.Sprintf("error /new_ballot : given tie-break %d is invalid or doesn't match #alts %d", req.TieBreak, req.Alts)
			w.Write([]byte(msg))
			return
		}
	}

	//Enregistre le nouveau ballot
	//Remarque : On pourrait mettre un mutex ici si on décidait d'utilisait 1 goroutine par ballot et de ne pas mettre de mutex sur la méthode vote
	var ballotId string = fmt.Sprintf("scrutin%d", rsa.countBallot)
	rsa.countBallot++
	rsa.ballotsList[ballotId], err = restagent.NewBallot(ballotId, req.Rule, req.Deadline, req.VoterIds, req.Alts, req.TieBreak)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		msg := fmt.Sprintf("error /new_ballot : can't create ballot %s. "+err.Error(), ballotId)
		w.Write([]byte(msg))
		return
	}
	var resp restagent.ResponseNewBallot = restagent.ResponseNewBallot{BallotId: ballotId}

	serial, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		msg := fmt.Sprint("error /new_ballot  : sérialisation de la réponse :", err.Error())
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(http.StatusCreated) //201
	w.Write(serial)
	fmt.Println("Liste ballots : ", rsa.ballotsList)
}
