package restclientagent

import (
	"log"
	"math/rand"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
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
