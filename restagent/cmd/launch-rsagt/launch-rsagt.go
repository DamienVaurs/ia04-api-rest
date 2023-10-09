package main

import (
	"fmt"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restserveragent"
)

func main() {
	server := restserveragent.NewRestServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
