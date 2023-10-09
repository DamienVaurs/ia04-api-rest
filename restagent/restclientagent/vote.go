package restclientagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/endpoints"
)

// Fonction qui réalise l'appel à l'API REST pour voter :
// http://localhost:8080/vote

func (rca *RestClientAgent) treatResponseVote(r *http.Response) error {
	//TODO : à voir si vraiment utile
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp restagent.Response
	err := json.Unmarshal(buf.Bytes(), &resp)

	return err
}

func (rca *RestClientAgent) doRequestVote(voteId string) (err error) {

	//Préparation de la requête
	req := restagent.RequestVote{
		AgentId: rca.id,
		VoteId:  voteId,
		Prefs:   rca.prefs,
		//TODO : voir si on ajoute option ou pas
	}

	// sérialisation de la requête
	url := rca.url + endpoints.Vote
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error by %s in /vote while marshalling request: %s", rca.id, err.Error())
	}

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		return fmt.Errorf("error by %s in /vote while sending request: %s", rca.id, err.Error())
	}
	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
	}
	return nil
}
