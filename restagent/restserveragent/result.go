package restserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
)

// Fonctions qui réalisent l'appel à l'API REST pour obtenir le résultat du vote :
// http://localhost:8080/result

// Décode la requête
func (*RestServerAgent) decodeResultRequest(r *http.Request) (req restagent.RequestResult, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		fmt.Println("Erreur de décodage de la requête /result : ", err)
		return
	}
	return
}

func (rsa *RestServerAgent) doCalcResult(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()
	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	req, err := rsa.decodeResultRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Println("Serveur recoit : ", r.URL, req)
	//Vérifie que le ballot existe
	_, found := rsa.ballotsList[req.BallotId]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("error /result : ballot %s does not exist", req.BallotId)
		w.Write([]byte(msg))
		return
	}

	// calcule de la réponse

	scf, err2 := comsoc.MajoritySCF(rsa.ballotsMap[req.BallotId])
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Error while processing result (SCF)"
		w.Write([]byte(msg))
		return
	}
	if len(scf) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Error : SCF ends with tie. No winner"
		w.Write([]byte(msg))
		return
	}
	resp := restagent.ResponseResult{}
	if len(scf) == 1 {
		resp.Winner = scf[0]
	}
	w.WriteHeader(http.StatusOK)
	serial, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Error while processing result (JSON) " + err.Error()
		w.Write([]byte(msg))
		return
	}
	w.Write(serial)
}
