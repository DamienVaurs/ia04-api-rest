package main

import (
	"math/rand"
	"sync"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/endpoints"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

/**
* Cette commande lance un agent tenant le scrutin et un agent votant qui réalisent les commandes suivantes :
* - POST /new_ballot : crée un scrutin
* - POST /vote : vote pour le scrutin
* - POST /result : récupère le résultat du scrutin
* - Affiche le résultat.
* Cela n'a pas grand intérêt en soit si ce n'est de tester le fonctionnement de l'API REST et la synchronisation entre agents.
*
**/

const nbAlts = 5 //nombre d'alternatives dans les préférences
//Remarque, on a choisit ce nombre arbitrairement

func main() {

	//Création de la requête pour créer un nouveau scrutin
	reqNewBallot := restagent.RequestNewBallot{
		Rule:     restagent.Majority,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
		VoterIds: []string{"ag_vote"}, //un seul votant car on ne lance qu'un client
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
	}

	//Création de la requête pour voter
	intPref := rand.Perm(nbAlts)
	altPref := make([]comsoc.Alternative, nbAlts)
	for i := 0; i < nbAlts; i++ {
		//conversion en []Alternative
		altPref[i] = comsoc.Alternative(intPref[i] + 1)
	}
	reqVote := restagent.RequestVote{
		AgentId:  "ag_vote",
		BallotId: "",
		Prefs:    altPref,
		Options:  nil,
	}

	//Création des canaux cin et cout
	listChannelsIn := make([]chan []string, 2) //Communication vers les agents
	channelOut := make(chan string)            //Communication des agents vers la goroutine principale

	listChannelsIn[0] = make(chan []string)
	listChannelsIn[1] = make(chan []string)
	//Lancement des agents
	agScrutin := restclientagent.NewRestClientBallotAgent("ag_scrut", endpoints.ServerHost+endpoints.ServerPort, reqNewBallot, listChannelsIn[0], channelOut)
	agVote := restclientagent.NewRestClientVoteAgent("ag_vote", endpoints.ServerHost+endpoints.ServerPort, reqVote, listChannelsIn[1], channelOut)

	wg := sync.WaitGroup{}
	wg.Add(2)

	//Lancement des agents

	go func() {
		//Lancement du créateur de scrutin
		defer wg.Done()
		agScrutin.Start()
	}()

	go func() {
		//Lancement du votant
		defer wg.Done()
		agVote.Start()
	}()
	//Attend la réception d'un scrutin
	ballotId := <-channelOut
	//Envoie de la liste des scrutins (1 seul) au votant
	listChannelsIn[1] <- []string{ballotId}
	//Attend la réception d'un message anoncant la fin des votes
	<-channelOut
	//Envoie d'un message à l'agent scrutin pour lancer le calcul des résultats
	listChannelsIn[0] <- []string{""}
	wg.Wait()
}
