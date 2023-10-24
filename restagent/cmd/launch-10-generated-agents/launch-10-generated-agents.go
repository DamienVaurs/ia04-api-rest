package main

import (
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/instances"
)

/**
* Cette commande lance un serveur et une flotte de 10 agents votants pour 8 scrutins (1 par méthode de vote et 2 témoins).
* Leurs préférences (et leur seuil) sont générées aléatoirement.
**/

func main() {
	instances.LaunchAgents(len(restagent.Rules)+2, 10, 5, instances.Init10VotingAgents)
}
