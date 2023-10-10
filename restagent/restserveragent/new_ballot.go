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

// Types de ballot autorisés
var typeBallot []string = []string{"majority", "approval", "condorcet", "copeland", "borda", "stv"}

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

func (rsa *RestServerAgent) doCreateNewBallot(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
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
	//Vérifie que le ballot n'existe pas déjà
	_, found := rsa.ballotsList[req.BallotId]
	if found {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("error /new_ballot : ballot %s already exists", req.BallotId)
		w.Write([]byte(msg))
		return
	}

	//Vérifie que le format de date est bon
	_, err = time.Parse("Mon Jan 02 15:04:05 MST 2006", req.Deadline)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("error /new_ballot : deadline %s is not in the right format", req.Deadline)
		w.Write([]byte(msg))
		return
	}

	//Vérifie que le type de ballot est autorisé
	var authorized = false
	for _, v := range typeBallot {
		if v == req.Rule {
			authorized = true
			break
		}
	}
	if !authorized {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("error : type %s is not authorized for ballot %s", req.Rule, req.BallotId)
		w.Write([]byte(msg))
		return
	}

	//Vérifie que les alternatives sont cohérentes avec le tie-break
	if req.Alts < 1 {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("error /new_ballot : number of alternatives %d is not correct", req.Alts)
		w.Write([]byte(msg))
		return
	}
	if req.TieBreak != nil {
		if len(req.TieBreak) != req.Alts {
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /new_ballot : number of alternatives %d is not correct with tie-break %v for ballot %s", req.Alts, req.TieBreak, req.BallotId)
			w.Write([]byte(msg))
			return
		}
		list := make([]comsoc.Alternative, len(req.TieBreak))
		copy(list, req.TieBreak)
		sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
		for i := 0; i < len(req.TieBreak)-1; i++ {
			if list[i]+1 != list[i+1] {
				w.WriteHeader(http.StatusBadRequest)
				msg := fmt.Sprintf("error /new_ballot : tie-break %v is not correct for ballot %s", req.TieBreak, req.BallotId)
				w.Write([]byte(msg))
				return
			}
		}
	}

	//Enregistre le nouveau ballot
	rsa.ballotsList[req.BallotId], err = restagent.NewBallot(req.BallotId, req.Rule, req.Deadline, req.VoterIds, req.Alts, req.TieBreak)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("error /new_ballot : can't create ballot %s. "+err.Error(), req.BallotId)
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(http.StatusOK)
	serial, err := json.Marshal(req)
	if err != nil {
		msg := fmt.Sprint("Erreur de sérialisation de la réponse /new_ballot : ", err)
		w.Write([]byte(msg))
		return
	}
	w.Write(serial)
	fmt.Println("Liste ballots : ", rsa.ballotsList)
}
