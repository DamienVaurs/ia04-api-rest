package restclientagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/endpoints"
)

// Fonctions qui réalisent l'appel à l'API REST pour obtenir le résultat du vote :
// http://localhost:8080/result

func (rca *RestClientAgent) treatResponseResults(r *http.Response) (resp restagent.ResponseResult, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	err = json.Unmarshal(buf.Bytes(), &resp)

	return
}

func (rca *RestClientAgent) doRequestResults(ballotId string) (res restagent.ResponseResult, err error) {

	// sérialisation de la requête
	url := rca.url + endpoints.Results

	//Création de la requête
	req := restagent.RequestResult{
		BallotId: ballotId,
	}

	// envoi de la requête
	data, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("error by %s in /request while marshalling request: %s", rca.id, err.Error())
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		return res, fmt.Errorf("error by %s in /request while sending request: %s", rca.id, err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	res, err = rca.treatResponseResults(resp)
	if err != nil {
		return res, fmt.Errorf("error by %s in /request while treating response: %s", rca.id, err.Error())
	}

	return
}
