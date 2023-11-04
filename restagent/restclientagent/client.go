package restclientagent

import (
	"fmt"
	"log"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
)

/******************  Corps d'un agent client ******************/
type RestClientAgentBase struct {
	Id   string          //id de l'agent
	url  string          //url du serveur
	cin  <-chan []string //channel pour recevoir la liste des scrutins
	cout chan<- string   //channel pour communiquer le nom de son scrutin
}

/******************  Agent créant un scrutin et calculant le résultat ******************/
type RestClientBallotAgent struct {
	RestClientAgentBase                            //attributs de base d'un agent client
	ReqNewBallot        restagent.RequestNewBallot //requête pour créer un nouveau scrutin
}

// Constructeur d'un agent créant un scrutin
func NewRestClientBallotAgent(id string, url string, reqNewBallot restagent.RequestNewBallot, cin <-chan []string, cout chan<- string) *RestClientBallotAgent {
	return &RestClientBallotAgent{
		RestClientAgentBase{id, url, cin, cout},
		reqNewBallot,
	}
}

// Méthode principale de l'Agent créant un scrutin
func (rcba *RestClientBallotAgent) Start() {
	//log.Printf("démarrage de l'agent pour la création d'un scrutin %s...", rcba.id)

	// Etape 1: Création du scrutin
	createdBallot, err := rcba.doRequestNewBallot(rcba.ReqNewBallot)
	if err != nil {
		log.Printf(rcba.Id, " error: ", err.Error()) //Remarque : on ne fait pas appel à log.Fatal car on veut que l'agent continue de fonctionner pour réaliser ses tâches
	} else {
		//log.Printf("/new_Ballot par [%s] créé avec succes : %s\n", rcba.id, createdBallot.BallotId)
	}
	// Etape 2: Envoie son scrutin à la goRoutine principale
	rcba.cout <- createdBallot.BallotId

	// Etape 3: Attente de la fin des votes de tous les agents, signalé par la goRoutine principale
	<-rcba.cin

	time.Sleep(6 * time.Second)

	// Etape 4: Récupération du résultat de chaque scrutin
	res, err := rcba.doRequestResults(createdBallot.BallotId)
	if err != nil {
		log.Printf(rcba.Id, "error: ", err.Error()) //Remarque : on ne fait pas appel à log.Fatal car on veut que l'agent continue de fonctionner pour réaliser ses tâches
	} else {
		Affichage(createdBallot.BallotId, rcba.ReqNewBallot.Rule, len(rcba.ReqNewBallot.VoterIds), res)
	}
}

/****************** Agent votant ******************/
type RestClientVoteAgent struct {
	RestClientAgentBase                       //attributs de base d'un agent client
	ReqVote             restagent.RequestVote //requête pour voter

}

// Constructeur d'un agent votant
func NewRestClientVoteAgent(id string, url string, reqVote restagent.RequestVote, cin <-chan []string, cout chan<- string) *RestClientVoteAgent {
	return &RestClientVoteAgent{
		RestClientAgentBase{id, url, cin, cout},
		reqVote,
	}
}

// Méthode principale de l'Agent votant
func (rcva *RestClientVoteAgent) Start() {

	//Etape 1 : Attente de la réception de l'ensemble des scrutins envoyés par la goRoutine principale
	listBallots := <-rcva.cin

	//Etape 2: Vote dans chaque scrutin

	for _, ballot := range listBallots {
		rcva.ReqVote.BallotId = ballot //Met l'Id du scrutin dans la requête
		err := rcva.doRequestVote(rcva.ReqVote)
		if err != nil {
			log.Printf(rcva.Id, " error: ", err.Error()) //Remarque : on ne fait pas appel à log.Fatal car on veut que l'agent continue de fonctionner pour réaliser ses tâches
		}
	}
	//Etape 3 : Envoie un message à la goRoutine principale pour signifier qu'il a terminé
	rcva.cout <- "fin"
}

func Affichage(id string, rule string, nbVoters int, res restagent.ResponseResult) {
	if rule != "condorcet" {
		fmt.Printf("=============================== RÉSULTATS POUR LE SCRUTIN %s ===============================\nTYPE DE SCRUTIN : %s\nNOMBRE DE VOTANTS : %d\nGAGNANT : %d\nCLASSEMENT : %v\n",
			id, rule, nbVoters, res.Winner, res.Ranking)
	} else {
		fmt.Printf("=============================== RÉSULTATS POUR LE SCRUTIN %s ===============================\nTYPE DE SCRUTIN : %s\nNOMBRE DE VOTANTS : %d\nGAGNANT : %d\n",
			id, rule, nbVoters, res.Winner)
	}
}
