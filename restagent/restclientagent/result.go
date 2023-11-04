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

func (rca *RestClientBallotAgent) treatResponseResults(r *http.Response) (resp restagent.ResponseResult, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	err = json.Unmarshal(buf.Bytes(), &resp)

	return
}

func (rca *RestClientBallotAgent) doRequestResults(ballotId string) (res restagent.ResponseResult, err error) {

	// sérialisation de la requête
	url := rca.url + endpoints.Results

	//Création de la requête
	req := restagent.RequestResult{
		BallotId: ballotId,
	}

	// envoi de la requête
	data, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("/result. error by %s in /result while marshalling request: %s", rca.Id, err.Error())
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement
	if err != nil {
		return res, fmt.Errorf("/result.error by %s in /result while sending request: %s", rca.Id, err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	res, err = rca.treatResponseResults(resp)
	if err != nil {
		return res, fmt.Errorf("/result.error by %s in /result while treating response: %s", rca.Id, err.Error())
	}

	return
}
