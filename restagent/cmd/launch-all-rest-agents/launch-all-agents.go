package main

import (
	"fmt"
	"log"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restserveragent"
)

func main() {
	const n = 10
	const url1 = ":8080"
	const url2 = "http://localhost:8080"

	clAgts := make([]restclientagent.RestClientAgent, 0, n)
	servAgt := restserveragent.NewRestServerAgent(url1)

	log.Println("démarrage du serveur...")
	go servAgt.Start()

	log.Println("démarrage des clients...")
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("id%02d", i)
		agt := restclientagent.NewRestClientAgent(id, url2, "vote") //TODO : décider si on sélectionne les préférences dans le Start() ou au dessus
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
		agt := restclientagent.NewRestClientAgent(id, url2, "results") //TODO : décider si on sélectionne les préférences dans le Start() ou au dessus
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
