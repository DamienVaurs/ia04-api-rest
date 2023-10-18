package instances

import (
	"log"
	"sync"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/endpoints"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restserveragent"
)

/**
* Abstraction des différents main lançant une flotte d'agents votants différents
*
* Cette fonction prend en entrée le paramètres :
* - n : nombre d'agents votants
* - nbAlts : nombre d'alternatives dans les préférences
* - generateAgentsFunc : fonction qui génère les agents votants
*
* Elle lance un serveur et une flotte d'agents votants, dont chacun réalise les commandes suivantes :
* - POST /new_ballot : crée un scrutin
* - Attend que les autres aient fini de créer leur scrutin
* - POST /vote : vote (ou tente de voter) pour chaque scrutin
* - Attend que les autres aient fini de voter pour chaque scrutin
* - POST /result : récupère le résultat de son scrutin
* - Affiche le résultat.
**/

func LaunchAgents(n int, nbAlts int, generateAgentsFunc func(longUrl string, n int, nbAlts int, listCin []chan []string, cout chan string) []restclientagent.RestClientAgent) {
	const url1 = endpoints.ServerPort
	const url2 = endpoints.ServerHost + endpoints.ServerPort
	servAgt := restserveragent.NewRestServerAgent(url1) //Serveur

	//Canaux de communication
	channelOut := make(chan string)           //channel par lequel les go routines communiquent à la goroutine principale
	channelInList := make([]chan []string, n) //contients l'ensemble des channels sur lesquels les agents vont lire (1 channel par agent)
	for i := 0; i < n; i++ {
		channelInList[i] = make(chan []string)
	}

	//Initialisation des agents votants (à la main pour choisir les paramètres)
	listAgents := generateAgentsFunc(url2, n, nbAlts, channelInList[:], channelOut)

	//Démarrage du serveur
	log.Println("démarrage du serveur...")
	go servAgt.Start()

	//Démarage des agents votants
	log.Println("démarrage des clients...")
	wg := sync.WaitGroup{}
	wg.Add(len(listAgents))

	for _, agt := range listAgents {
		go func(agt restclientagent.RestClientAgent) { //Passage par fonction lamda pour capturer la valeur de l'itération par la goroutine
			defer wg.Done()
			agt.Start()
		}(agt)
	}

	log.Print("\n\n============================= Création des scrutins =============================\n\n")
	listBallots := make([]string, len(listAgents))
	//Attend la réception des n scrutins
	for i := 0; i < len(listAgents); i++ {
		ballotId := <-channelOut
		listBallots[i] = ballotId
	}

	log.Print("\n\n============================= Votes des agents =============================\n\n")
	//Envoie de la liste des scrutins à tous les agents
	for _, channelIn := range channelInList {
		channelIn <- listBallots
	}

	//Attend la réception de n messages anoncant la fin des votes
	for i := 0; i < len(listAgents); i++ {
		<-channelOut
	}
	log.Print("\n\n============================= Réception des résultats =============================\n\n")
	//Envoie d'un message pour lancer l'appel aux résultats
	for _, channelIn := range channelInList {
		channelIn <- []string{"fin"}
	}
	wg.Wait()

	log.Print("\n\n============================= FIN PROGRAMME =============================\n\n")
}
