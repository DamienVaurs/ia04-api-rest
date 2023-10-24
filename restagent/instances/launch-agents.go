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
* Elle lance un serveur et une flotte d'agents scrutins et votants, dont chacun réalise les commandes suivantes :
* - POST /new_ballot : crée un scrutin
* - Attend que les autres aient fini de créer leur scrutin
* - POST /vote : vote (ou tente de voter) pour chaque scrutin
* - Attend que les autres aient fini de voter pour chaque scrutin
* - POST /result : récupère le résultat de son scrutin
* - Affiche le résultat.
**/

func LaunchAgents(nbBallot int, nbVotant int, nbAlts int, generateAgentsFunc func(longUrl string, n int, nbBallot int, nbAlts int, listCinVotants []chan []string, listCinBallots []chan []string, cout chan string) ([]restclientagent.RestClientVoteAgent, []restclientagent.RestClientBallotAgent)) {
	const url1 = endpoints.ServerPort
	const url2 = endpoints.ServerHost + endpoints.ServerPort
	servAgt := restserveragent.NewRestServerAgent(url1) //Serveur

	//Canaux de communication
	channelOut := make(chan string)                             //channel par lequel les go routines communiquent à la goroutine principale
	channelInListBallotAgent := make([]chan []string, nbBallot) //channel sur lesquels liront les BallotAgents
	channelInListVoteAgent := make([]chan []string, nbVotant)   //channel sur lesquels liront les VoteAgents
	for i := 0; i < nbBallot; i++ {
		channelInListBallotAgent[i] = make(chan []string)
	}
	for i := 0; i < nbVotant; i++ {
		channelInListVoteAgent[i] = make(chan []string)
	}

	//Initialisation des agents votants (à la main pour choisir les paramètres)
	listVoteAgents, listBallotAgents := generateAgentsFunc(url2, nbVotant, nbBallot, nbAlts, channelInListVoteAgent[:], channelInListBallotAgent[:], channelOut)

	//Démarrage du serveur
	log.Println("démarrage du serveur...")
	go servAgt.Start()

	//Démarage des agents clients
	log.Println("démarrage des clients...")

	log.Print("\n\n============================= Lancement des Agents et Création des scrutins =============================\n\n")

	wg := sync.WaitGroup{}
	wg.Add(len(listBallotAgents) + len(listVoteAgents))

	//Lancement des agents scrutins
	for _, agt := range listBallotAgents {
		go func(agt restclientagent.RestClientBallotAgent) { //Passage par fonction lamda pour capturer la valeur de l'itération par la goroutine
			defer wg.Done()
			agt.Start()
		}(agt)
	}

	//Lancement des agents votants
	for _, agt := range listVoteAgents {
		go func(agt restclientagent.RestClientVoteAgent) { //Passage par fonction lamda pour capturer la valeur de l'itération par la goroutine
			defer wg.Done()
			agt.Start()
		}(agt)
	}

	listBallots := make([]string, len(listBallotAgents))
	//Attend la réception des len(listBallotAgents) scrutins
	for i := 0; i < len(listBallotAgents); i++ {
		ballotId := <-channelOut
		listBallots[i] = ballotId
	}

	log.Print("\n\n============================= Votes des agents =============================\n\n")
	//Envoie de la liste des scrutins à tous les agents votants
	for _, channelIn := range channelInListVoteAgent {
		channelIn <- listBallots
	}

	//Attend la réception de n messages de la part des votants anoncant la fin des votes
	for i := 0; i < len(listVoteAgents); i++ {
		<-channelOut
	}
	log.Print("\n\n============================= Réception des résultats =============================\n\n")
	log.Print("Attente de la fin des scrutins...\n\n")
	//Envoie d'un message aux agents scrutins pour lancer le calcul des résultats
	for _, channelIn := range channelInListBallotAgent {
		channelIn <- []string{"fin"}
	}
	wg.Wait()

	log.Print("\n\n============================= FIN PROGRAMME =============================\n\n")
}
