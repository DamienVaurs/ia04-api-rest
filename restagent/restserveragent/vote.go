package restserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
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
	//Vérifie que le ballot existe
	_, found := rsa.ballotsList[req.VoteId]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("error /vote : ballot %s does not exist", req.VoteId)
		w.Write([]byte(msg))
		return
	}

	//TODO : vérifir que l'agent n'a pas déjà voté

	//Enregistre le vote pour le ballot
	rsa.ballotsMap[req.VoteId] = append(rsa.ballotsMap[req.VoteId], req.Prefs)

	w.WriteHeader(http.StatusOK)
	serial, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Write(serial)
}
