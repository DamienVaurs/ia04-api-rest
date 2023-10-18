package main

import "gitlab.utc.fr/milairhu/ia04-api-rest/restagent/instances"

/**
* Cette commande lance un serveur et une flotte d'agents votants
* pour calculer les résultats de l'exemple de la diapositive 34, chapitre 3
* du cours
**/

func main() {
	instances.LaunchAgents(3, 5+4+2+6+8+2, 4, instances.InitChap3Diap34)
}
