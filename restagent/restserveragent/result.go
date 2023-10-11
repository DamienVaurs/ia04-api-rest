package restserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	//Vérifie que la date de fin est passée
	if rsa.ballotsList[req.BallotId].Deadline.After(time.Now()) {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("error /result : ballot %s is not finished yet", req.BallotId)
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

		//TODO : appliquer le ranking et le tie-break pour Approval
	} else if rsa.ballotsList[req.BallotId].Rule == "condorcet" {
		//TODO : appliquer le ranking et le tie-break pour Approval
		scf, err := comsoc.CondorcetWinner(rsa.ballotsMap[req.BallotId])
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
		var scfVote func(comsoc.Profile) ([]comsoc.Alternative, error)
		var swfVote func(comsoc.Profile) (comsoc.Count, error)
		switch rsa.ballotsList[req.BallotId].Rule {
		case "borda":
			scfVote = comsoc.BordaSCF
			swfVote = comsoc.BordaSWF
		case "copeland":
			scfVote = comsoc.CopelandSCF
			swfVote = comsoc.CopelandSWF
		case "majority":
			scfVote = comsoc.MajoritySCF
			swfVote = comsoc.MajoritySWF
		case "stv":
			scfVote = comsoc.STV_SCF
			swfVote = comsoc.STV_SWF
		default:
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("error /result : type %s is not authorized for ballot %s", rsa.ballotsList[req.BallotId].Rule, req.BallotId)
			w.Write([]byte(msg))
			return
		}

		//Si on a un tie-break, on l'applique pour avoir le meilleur élément et le classement
		var tieBreak func([]comsoc.Alternative) (comsoc.Alternative, error)
		if rsa.ballotsList[req.BallotId].TieBreak != nil {
			tieBreak = comsoc.TieBreakFactory(rsa.ballotsList[req.BallotId].TieBreak)
			swfFunc := comsoc.SWFFactory(swfVote, tieBreak)
			res, err := swfFunc(rsa.ballotsMap[req.BallotId])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				msg := fmt.Sprintf("error /result : can't process SWF with Tie-break for ballot %s of type %s. "+err.Error(), req.BallotId, rsa.ballotsList[req.BallotId].Rule)
				w.Write([]byte(msg))
				return
			}
			resp.Winner = res[0]
			resp.Ranking = res
		} else {
			//Si on n'a pas de tie-break
			scf, err := scfVote(rsa.ballotsMap[req.BallotId])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				msg := fmt.Sprintf("error /result : can't process SCF for ballot %s of type %s. "+err.Error(), req.BallotId, rsa.ballotsList[req.BallotId].Rule)
				w.Write([]byte(msg))
				return
			}
			swf, err := swfVote(rsa.ballotsMap[req.BallotId])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				msg := fmt.Sprintf("error /result : can't process SWF for ballot %s of type %s. "+err.Error(), req.BallotId, rsa.ballotsList[req.BallotId].Rule)
				w.Write([]byte(msg))
				return
			}
			//Il faut calculer l'ordre du SWF
			ranking := comsoc.MakeRanking(swf)
			resp.Winner = scf[0]
			resp.Ranking = ranking
		}

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
