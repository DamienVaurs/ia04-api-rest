package instances

import (
	"fmt"
	"log"
	"sync"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
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

	//Initialisation  agents votants (à la main pour choisir les paramètres)
	listVoteAgents, listBallotAgents := generateAgentsFunc(url2, nbVotant, nbBallot, nbAlts, channelInListVoteAgent[:], channelInListBallotAgent[:], channelOut)

	//Démarrage du serveur
	log.Println("démarrage du serveur...")
	go servAgt.Start()

	//Démarage des agents clients
	log.Println("démarrage des clients...")

	fmt.Print("\n\n============================= Lancement des Agents et Création des scrutins =============================\n\n")

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

	fmt.Print("\n\n============================= Votes des agents =============================\n\n")
	//Envoie de la liste des scrutins à tous les agents votants
	for _, channelIn := range channelInListVoteAgent {
		channelIn <- listBallots
	}

	//Attend la réception de n messages de la part des votants anoncant la fin des votes
	for i := 0; i < len(listVoteAgents); i++ {
		<-channelOut
	}

	fmt.Printf("\nPROFILS DES VOTANTS AU(X) SCRUTIN(S)\n")

	for _, ballot := range listBallotAgents {
		fmt.Printf("\n\nScrutin %v\n", ballot.RestClientAgentBase.Id)

		matrice := make([][]int, nbAlts+1)                             // matrice des préférences des votants
		listProfiles := make([]restclientagent.RestClientVoteAgent, 0) // liste des profils déjà inclus dans la matrice
		for i, agent := range listVoteAgents {
			if canVote(ballot.ReqNewBallot.VoterIds, agent.Id) {
				k := 1 // nombre de votants ayant le même profil
				for j, agent2 := range listVoteAgents {
					// si l'agent j n'est pas l'agent i et que les préférences de l'agent i et de l'agent j sont égales
					if i != j && areSlicesEqual(agent.ReqVote.Prefs, agent2.ReqVote.Prefs) && canVote(ballot.ReqNewBallot.VoterIds, agent2.Id) {
						listProfiles = append(listProfiles, agent2) // on ajoute l'agent j à la liste des profils déjà inclus dans la matrice
						k++                                         // on incrémente le nombre de votants ayant le même profil
					}

				}
				// si l'agent i n'est pas déjà dans la matrice
				if notInSlice(agent, listProfiles) {
					matrice[0] = append(matrice[0], k)
					for p, alt := range agent.ReqVote.Prefs {
						matrice[p+1] = append(matrice[p+1], int(alt))
					}
				}
			}
		}

		for i := range matrice {
			fmt.Printf("%v\n", matrice[i])
			if i == 0 {
				for j := 0; j < len(matrice[0]); j++ {
					fmt.Printf("--")
				}
				fmt.Printf("-\n")
			}

		}
	}

	fmt.Print("\n\n============================= Réception des résultats =============================\n\n")
	//Envoie d'un message aux agents scrutins pour lancer le calcul des résultats
	for _, channelIn := range channelInListBallotAgent {
		channelIn <- []string{"fin"}
	}
	wg.Wait()

	fmt.Print("\n\n============================= FIN PROGRAMME =============================\n\n")
}

// Fonction qui compare deux slices d'Alternatives pour déterminer si elles sont égales
func areSlicesEqual(a, b []comsoc.Alternative) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Fonction qui compare un agent avec une liste d'agents pour déterminer si l'agent n'est pas dans la liste
func notInSlice(a restclientagent.RestClientVoteAgent, b []restclientagent.RestClientVoteAgent) bool {
	for _, v := range b {
		if v.ReqVote.AgentId == a.ReqVote.AgentId {
			return false
		}
	}
	return true
}

// Fonction qui renvoie True si un agent a le droit de voter à un scrutin
func canVote(voterIDs []string, agentId string) bool {
	for _, v := range voterIDs {
		if v == agentId {
			return true
		}
	}
	return false
}
