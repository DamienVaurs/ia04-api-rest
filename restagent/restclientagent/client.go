package restclientagent

import (
	"fmt"
	"log"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
)

type RestClientAgent struct {
	id           string                     //id de l'agent
	url          string                     //url du serveur
	reqNewBallot restagent.RequestNewBallot //requête pour créer un nouveau scrutin
	reqVote      restagent.RequestVote      //requête pour voter
	cin          <-chan []string            //channel pour recevoir la liste des scrutins
	cout         chan<- string              //channel pour communiquer le nom de son scrutin

}

// Constructeur d'un agent votant
func NewRestClientAgent(id string, url string, reqNewBallot restagent.RequestNewBallot, reqVote restagent.RequestVote, cin <-chan []string, cout chan<- string) *RestClientAgent {

	return &RestClientAgent{id, url, reqNewBallot, reqVote, cin, cout}
}

// Méthode principale de l'Agent
func (rca *RestClientAgent) Start() {
	log.Printf("démarrage de %s...", rca.id)

	//Etape 1: Création du scrutin
	createdBallot, err := rca.doRequestNewBallot(rca.reqNewBallot)
	if err != nil {
		log.Printf(rca.id, " error: ", err.Error()) //Remarque : on ne fait pas appel à log.Fatal car on veut que l'agent continue de fonctionner pour réaliser ses tâches
	} else {
		log.Printf("/new_Ballot par [%s] créé avec succes : %s\n", rca.id, createdBallot.BallotId)
	}
	//Envoie son scrutin à la goRoutine principale
	rca.cout <- createdBallot.BallotId

	//Attente de la réception de l'ensemble des scrutins envoyés par la goRoutine principale
	fmt.Printf("[%s]: attente des scrutins...\n", rca.id)
	listBallots := <-rca.cin

	//Etape 2: Vote dans chaque scrutin

	for _, ballot := range listBallots {
		rca.reqVote.BallotId = ballot //Met l'Id du scrutin dans la requête
		err := rca.doRequestVote(rca.reqVote)
		if err != nil {
			log.Printf(rca.id, " error: ", err.Error()) //Remarque : on ne fait pas appel à log.Fatal car on veut que l'agent continue de fonctionner pour réaliser ses tâches
		} else {
			log.Printf("/vote par [%s] envoyé avec succes : %d\n", rca.id, rca.reqVote.Prefs)
		}
	}
	//Envoie un message à la goRoutine principale pour signifier qu'il a terminé
	rca.cout <- "fin"

	//Attente de la fin des votes de tous les agents, signalé par la goRoutine principale
	fmt.Printf("[%s] : attente de la fin des votes...\n", rca.id)
	<-rca.cin

	//Etape 3: Récupération du résultat de chaque scrutin
	res, err := rca.doRequestResults(createdBallot.BallotId)
	if err != nil {
		log.Printf(rca.id, "error: ", err.Error()) //Remarque : on ne fait pas appel à log.Fatal car on veut que l'agent continue de fonctionner pour réaliser ses tâches
	} else {
		log.Printf("[%s] : resultat recu pour le scrutin %s de type %s = %d\n", rca.id, createdBallot, rca.reqNewBallot.Rule, res)
	}
}
