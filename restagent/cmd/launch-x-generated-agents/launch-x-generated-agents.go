package main

import (
	"fmt"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/instances"
)

func main() {

	var nbAgents int
	var nbBallot int
	var nbAlts int

	fmt.Println("Combien d'agents votants?")
	fmt.Scanln(&nbAgents)
	fmt.Println("Combien de scrutins ?")
	fmt.Scanln(&nbBallot)
	fmt.Println("Combien d'alternatives ?")
	fmt.Scanln(&nbAlts)

	instances.LaunchAgents(nbBallot, nbAgents, nbAlts, instances.InitVotingAgents)
}
