package main

import "gitlab.utc.fr/milairhu/ia04-api-rest/restagent/instances"

/**
* Cette commande lance un serveur et une flotte d'agents votants
* pour calculer les résultats de l'exemple de la diapositive 34, chapitre 2
* du cours
**/

func main() {
	instances.LaunchAgents(1, 5+4+3, 4, instances.InitChap3Diap36)
}
