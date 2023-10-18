package main

import (
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/instances"
)

/**
* Cette commande lance un serveur et une flotte de 10 agents votants.
* Leurs préférences sont générées aléatoirement.
* On a essaye de lancer suffisament de cas différents pour tester
* l'ensemble de l'application
**/

func main() {
	instances.LaunchAgents(len(restagent.Rules), 10, 5, instances.Init10VotingAgents)
}
