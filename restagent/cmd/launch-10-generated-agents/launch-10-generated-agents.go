package main

import (
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/instances"
)

/**
* Cette commande lance un serveur et une flotte de 10 agents votants pour 6 scrutins (1 par méthode de vote).
* Leurs préférences (et leur seuil) sont générées aléatoirement.
**/

func main() {
	instances.LaunchAgents(len(restagent.Rules), 10, 5, instances.Init10VotingAgents)
}
