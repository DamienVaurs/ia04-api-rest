package main

import (
	"fmt"
	"log"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/endpoints"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restserveragent"
)

func main() {
	const n = 10
	const url1 = endpoints.ServerPort
	const url2 = endpoints.ServerHost + endpoints.ServerPort

	clAgts := make([]restclientagent.RestClientAgent, 0, n)
	servAgt := restserveragent.NewRestServerAgent(url1)

	log.Println("démarrage du serveur...")
	go servAgt.Start()

	log.Println("démarrage des clients...")
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("id%02d", i)

		req := restagent.RequestNewBallot{
			Rule:     "majority",
			Deadline: 0,
			VoterIds: []string{"id00", "id01", "id02", "id03", "id04", "id05", "id06", "id07", "id08", "id09"},
			Alts:     []string{"alt00", "alt01", "alt02", "alt03", "alt04"},
			TieBreak: "random",
		}
		agt := restclientagent.NewRestClientAgent(id, url2, "vote", req) //TODO : décider si on sélectionne les préférences dans le Start() ou au dessus
		clAgts = append(clAgts, *agt)
	}

	for _, agt := range clAgts {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt restclientagent.RestClientAgent) {
			go agt.Start()
		}(agt)
	}

	//Lance 3 agents pour avoir les résultats
	for i := 0; i < 3; i++ {
		id := fmt.Sprintf("idRes%02d", i)
		agt := restclientagent.NewRestClientAgent(id, url2, "result") //TODO : décider si on sélectionne les préférences dans le Start() ou au dessus
		clAgts = append(clAgts, *agt)
	}
	clAgtsRes := make([]restclientagent.RestClientAgent, 0, n)
	for _, agt := range clAgtsRes {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt restclientagent.RestClientAgent) {
			go agt.Start()
		}(agt)
	}
	fmt.Scanln()
}
