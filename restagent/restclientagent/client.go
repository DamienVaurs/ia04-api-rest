package restclientagent

import (
	"fmt"
	"log"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
)

/******************  Corps d'un agent client ******************/
type RestClientAgentBase struct {
	id   string          //id de l'agent
	url  string          //url du serveur
	cin  <-chan []string //channel pour recevoir la liste des scrutins
	cout chan<- string   //channel pour communiquer le nom de son scrutin
}

/******************  Agent créant un scrutin et calculant le résultat ******************/
type RestClientBallotAgent struct {
	RestClientAgentBase                            //attributs de base d'un agent client
	reqNewBallot        restagent.RequestNewBallot //requête pour créer un nouveau scrutin
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
	log.Printf("démarrage de l'agent pour la création d'un scrutin %s...", rcba.id)

	// Etape 1: Création du scrutin
	createdBallot, err := rcba.doRequestNewBallot(rcba.reqNewBallot)
	if err != nil {
		log.Printf(rcba.id, " error: ", err.Error()) //Remarque : on ne fait pas appel à log.Fatal car on veut que l'agent continue de fonctionner pour réaliser ses tâches
	} else {
		log.Printf("/new_Ballot par [%s] créé avec succes : %s\n", rcba.id, createdBallot.BallotId)
	}
	// Etape 2: Envoie son scrutin à la goRoutine principale
	rcba.cout <- createdBallot.BallotId

	// Etape 3: Attente de la fin des votes de tous les agents, signalé par la goRoutine principale
	fmt.Printf("[%s] : attente de la fin des votes...\n", rcba.id)
	<-rcba.cin

	// Etape 4: Récupération du résultat de chaque scrutin
	res, err := rcba.doRequestResults(createdBallot.BallotId)
	if err != nil {
		log.Printf(rcba.id, "error: ", err.Error()) //Remarque : on ne fait pas appel à log.Fatal car on veut que l'agent continue de fonctionner pour réaliser ses tâches
	} else {
		log.Printf("[%s] : resultat recu pour le scrutin %s de type %s = %d\n", rcba.id, createdBallot, rcba.reqNewBallot.Rule, res)
	}
}

/****************** Agent votant ******************/
type RestClientVoteAgent struct {
	RestClientAgentBase                       //attributs de base d'un agent client
	reqVote             restagent.RequestVote //requête pour voter

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
	log.Printf("démarrage de l'agent pour le vote %s...", rcva.id)

	//AEtape 1 : Attente de la réception de l'ensemble des scrutins envoyés par la goRoutine principale
	fmt.Printf("[%s]: attente des scrutins...\n", rcva.id)
	listBallots := <-rcva.cin

	//Etape 2: Vote dans chaque scrutin

	for _, ballot := range listBallots {
		rcva.reqVote.BallotId = ballot //Met l'Id du scrutin dans la requête
		err := rcva.doRequestVote(rcva.reqVote)
		if err != nil {
			log.Printf(rcva.id, " error: ", err.Error()) //Remarque : on ne fait pas appel à log.Fatal car on veut que l'agent continue de fonctionner pour réaliser ses tâches
		} else {
			log.Printf("/vote par [%s] envoyé avec succes : %d\n", rcva.id, rcva.reqVote.Prefs)
		}
	}
	//Etape 3 : Envoie un message à la goRoutine principale pour signifier qu'il a terminé
	rcva.cout <- "fin"
}
