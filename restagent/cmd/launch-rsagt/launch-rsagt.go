package main

import (
	"fmt"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/endpoints"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restserveragent"
)

func main() {
	server := restserveragent.NewRestServerAgent(endpoints.ServerPort)
	server.Start()
	fmt.Scanln()
}
