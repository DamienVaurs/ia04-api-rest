package main

import (
	"math/rand"
	"sync"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/endpoints"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

/**
* Cette commande lance un agent votant qui réalise les commandes suivantes :
* - POST /new_ballot : crée un scrutin
* - POST /vote : vote pour le scrutin
* - POST /result : récupère le résultat du scrutin
* - Affiche le résultat.
* Cela n'a pas grand intérêt en soit si ce n'est de tester le fonctionnement de l'API REST.
* Pour un exemple plus intéressant, voir la commande launch-all-agent qui lance une flotte d'agents.
*
**/

const nbAlts = 5 //nombre d'alternatives dans les préférences
//Remarque, on a choisit ce nombre arbitrairement

func main() {

	//Création de la requête pour créer un nouveau scrutin
	reqNewBallot := restagent.RequestNewBallot{
		Rule:     restagent.Majority,
		Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
		VoterIds: []string{"id1"},        //un seul votant car on ne lance qu'un client
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
		AgentId:  "id1",
		BallotId: "",
		Prefs:    altPref,
		Options:  nil,
	}

	//Création des canaux cin et cout
	channelIn := make(chan []string)
	channelOut := make(chan string)

	//Lancement de l'agent
	ag := restclientagent.NewRestClientAgent("id1", endpoints.ServerHost+endpoints.ServerPort, reqNewBallot, reqVote, channelIn, channelOut)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		ag.Start()
	}()
	//Attend la réception d'un scrutin
	ballotId := <-channelOut
	//Envoie de la liste des scrutins (1 seul)
	channelIn <- []string{ballotId}
	//Attend la réception d'un message anoncant la fin des votes
	<-channelOut
	//Envoie d'un message pour lancer l'appel aux résultats
	channelIn <- []string{""}
	wg.Wait()
}
