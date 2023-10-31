package main

import "gitlab.utc.fr/milairhu/ia04-api-rest/restagent/instances"

/**
* Cette commande lance un serveur et une flotte d'agents votants
* pour calculer les résultats de 2 scrutins STV, un impiquant l'usage de Tie-Break, l'autre non.
**/

func main() {
	instances.LaunchAgents(2, 5+4+3, 4, instances.InitSTV)
}
