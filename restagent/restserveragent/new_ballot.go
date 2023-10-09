package restserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
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
	_, found := rsa.ballotsMap[req.BallotId]
	if found {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("error /new_ballot : ballot %s already exists", req.BallotId)
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
	//Enregistre le nouveau ballot
	rsa.ballotsList[req.BallotId] = restagent.Ballot{
		BallotId: req.BallotId,
		Rule:     req.Rule,
		Deadline: req.Deadline,
		VoterIds: req.VoterIds,
		Alts:     req.Alts,
		TieBreak: req.TieBreak}

	w.WriteHeader(http.StatusOK)
	serial, err := json.Marshal(req)
	if err != nil {
		msg := fmt.Sprint("Erreur de sérialisation de la réponse /new_ballot : ", err)
		w.Write([]byte(msg))
		return
	}
	w.Write(serial)
}
