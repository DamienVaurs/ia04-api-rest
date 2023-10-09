package main

import (
	"fmt"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

func main() {
	ag := restclientagent.NewRestClientAgent("id1", "http://localhost:8080", "vote") //TODO : décider si on sélectionne les préférences dans le Start() ou au dessus
	ag.Start()
	fmt.Scanln()
}
