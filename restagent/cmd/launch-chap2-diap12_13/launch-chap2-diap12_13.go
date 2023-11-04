package main

import "gitlab.utc.fr/milairhu/ia04-api-rest/restagent/instances"

/**
* Cette commande lance un serveur et une flotte d'agents votants
* pour calculer les résultats de l'exemple des diapos 12 et 13, chapitre 2
* du cours
**/

func main() {
	instances.LaunchAgents(2, 6, 3, instances.InitChap3Diap12_13)
}
