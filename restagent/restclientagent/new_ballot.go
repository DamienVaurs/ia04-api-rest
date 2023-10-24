package restclientagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/endpoints"
)

// Fonctions qui réalisent l'appel à l'API REST pour voter :
// http://localhost:8080/new_ballot

func (rca *RestClientBallotAgent) treatResponseNewBallot(r *http.Response) (resp restagent.ResponseNewBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	err = json.Unmarshal(buf.Bytes(), &resp)

	return
}

func (rca *RestClientBallotAgent) doRequestNewBallot(req restagent.RequestNewBallot) (res restagent.ResponseNewBallot, err error) {

	// sérialisation de la requête
	url := rca.url + endpoints.NewBallot
	data, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("/new_ballot. error by %s in /new_ballot while marshalling request: %s", rca.id, err.Error())
	}

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		return res, fmt.Errorf("/new_ballot. error by %s in /new_ballot while sending request: %s", rca.id, err.Error())
	}
	if resp.StatusCode != http.StatusCreated {

		return res, fmt.Errorf("/new_ballot. [%d] %s", resp.StatusCode, resp.Status)
	}
	res, err = rca.treatResponseNewBallot(resp)
	if err != nil {
		return res, fmt.Errorf("/new_ballot. error by %s in /new_ballot while treating response: %s", rca.id, err.Error())
	}

	return
}
