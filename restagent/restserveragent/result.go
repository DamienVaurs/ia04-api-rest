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

	// calcule de la réponse en fonction de Ballot.Rule
	resp := restagent.ResponseResult{}
	if rsa.ballotsList[req.BallotId].Rule == "approval" {
		//Vérifie que le ballot a bien un seuil
		//TODO : appliquer le threshold pour ApprovalSCF
		scf, err := comsoc.ApprovalSCF(rsa.ballotsMap[req.BallotId], []int{1})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /result : can't process SCF for ballot %s of type %s", req.BallotId, rsa.ballotsList[req.BallotId].Rule)
			w.Write([]byte(msg))
			return
		}
		resp.Winner = scf[0]
		w.WriteHeader(http.StatusOK)
		serial, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /result : can't serialize response for ballot %s of type %s", req.BallotId, rsa.ballotsList[req.BallotId].Rule)
			w.Write([]byte(msg))
			return
		}
		w.Write(serial)
	} else {
		var methodVote func(comsoc.Profile) ([]comsoc.Alternative, error)
		switch rsa.ballotsList[req.BallotId].Rule {
		case "borda":
			methodVote = comsoc.BordaSCF
		case "condorcet":
			methodVote = comsoc.CondorcetWinner
		case "copeland":
			methodVote = comsoc.CopelandSCF
		case "majority":
			methodVote = comsoc.MajoritySCF
		case "stv":
			methodVote = comsoc.STV_SCF
		default:
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /result : type %s is not authorized for ballot %s", rsa.ballotsList[req.BallotId].Rule, req.BallotId)
			w.Write([]byte(msg))
			return
		}
		fmt.Println(rsa.ballotsMap[req.BallotId])
		scf, err := methodVote(rsa.ballotsMap[req.BallotId])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /result : can't process SCF for ballot %s of type %s. "+err.Error(), req.BallotId, rsa.ballotsList[req.BallotId].Rule)
			w.Write([]byte(msg))
			return
		}
		resp.Winner = scf[0]
		w.WriteHeader(http.StatusOK)
		serial, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /result : can't serialize response for ballot %s of type %s", req.BallotId, rsa.ballotsList[req.BallotId].Rule)
			w.Write([]byte(msg))
			return
		}
		w.Write(serial)
	}

}
