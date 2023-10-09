package restclientagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/endpoints"
)

const nbAlt = 5 //nombre d'alternatives dans les préférences

type RestClientAgent struct {
	id     string
	url    string
	prefs  []comsoc.Alternative
	action string
}

func NewRestClientAgent(id string, url string, action string) *RestClientAgent {

	src := rand.Perm(nbAlt) //TODO vérifier que ça fait bien ça
	dest := make([]comsoc.Alternative, nbAlt)
	for i, v := range src {
		dest[i] = comsoc.Alternative(v)
	}
	return &RestClientAgent{id, url, dest, action}
}

func (rca *RestClientAgent) doRequestVote() (res []comsoc.Alternative, err error) {
	req := restagent.Request{
		//Mettre les champs de la requete
		Preferences: rca.prefs,
	}

	// sérialisation de la requête
	url := rca.url + endpoints.Vote
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		return []comsoc.Alternative{}, err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return []comsoc.Alternative{}, err
	}
	res = rca.treatResponseVote(resp)
	return res, nil
}

func (rca *RestClientAgent) doRequestResults() (res comsoc.Alternative, err error) {

	// sérialisation de la requête
	url := rca.url + endpoints.Results

	// envoi de la requête
	resp, err := http.Get(url)

	// traitement de la réponse
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	res, err = rca.treatResponseResults(resp)
	if err != nil {
		return -1, err
	}

	return
}

func (rca *RestClientAgent) treatResponseVote(r *http.Response) []comsoc.Alternative {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp restagent.Response
	json.Unmarshal(buf.Bytes(), &resp)

	return resp.Result
}

func (rca *RestClientAgent) treatResponseResults(r *http.Response) (comsoc.Alternative, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp restagent.Response
	json.Unmarshal(buf.Bytes(), &resp)

	if len(resp.Result) < 1 {
		return -1, fmt.Errorf("error : scf ends with tie. No winner")
	} else if len(resp.Result) > 1 {
		return -1, fmt.Errorf("error : result contains more than 1 element")
	}

	return resp.Result[0], nil
}

func (rca *RestClientAgent) Start() {
	log.Printf("démarrage de %s", rca.id)
	if rca.action == "vote" {
		res, err := rca.doRequestVote()

		if err != nil {
			log.Fatal(rca.id, " error:", err.Error())
		} else {
			log.Printf("Vote [%s] = %d\n", rca.id, res)
		}
	} else if rca.action == "results" {
		res, err := rca.doRequestResults()
		if err != nil {
			log.Fatal(rca.id, "error:", err.Error())
		} else {
			log.Printf("Resultat [%s] = %d\n", rca.id, res)
		}
	}
}
