package restagent

import (
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
)

// Types utilisés pour la requête /new_ballot
type Ballot struct {
	BallotId   string               //Identifiant du ballot
	Rule       string               //Méthode de vote
	Deadline   time.Time            //Date limite de vote
	VoterIds   []string             //Liste des agents pouvant voter
	Alts       int                  //Alternatives de 1 à Alts
	TieBreak   []comsoc.Alternative //Ordre de préférence des alternatives en cas d'égalité
	HaveVoted  []string             //Noms des agents ayant voté
	Thresholds map[string]int       //Contient les seuils de chaque votant (pour vote par approbation)
}

// Constructeur d'un Ballot
func NewBallot(ballotId string, rule string, deadline string, voterIds []string, alts int, tieBreak []comsoc.Alternative) (Ballot, error) {
	//Vérifie que le format de date est bon
	date, err := time.Parse(time.RFC3339, deadline)
	if err != nil {
		return Ballot{}, err
	}
	l := make([]string, len(voterIds))
	t := make(map[string]int)
	return Ballot{
		BallotId:   ballotId,
		Rule:       rule,
		Deadline:   date,
		VoterIds:   voterIds,
		Alts:       alts,
		TieBreak:   tieBreak,
		HaveVoted:  l,
		Thresholds: t}, nil
}

type RequestNewBallot struct {
	Rule     string               `json:"rule"`      //Méthode de vote
	Deadline string               `json:"deadline"`  //Date limite de vote
	VoterIds []string             `json:"voter-ids"` //Liste des agents pouvant voter
	Alts     int                  `json:"#alts"`     //Alternatives de 1 à Alts
	TieBreak []comsoc.Alternative `json:"tie-break"` //Ordre de préférence des alternatives en cas d'égalité
}

type ResponseNewBallot struct {
	//Objet renvoyé si code 201
	BallotId string `json:"ballot-id"` //Id du scrutin créé
}

// Type utilisé pour la requête /vote
type RequestVote struct {
	AgentId  string               `json:"agent-id"`  //Id de l'agent votant
	BallotId string               `json:"ballot-id"` //Id du scrutin auquel on vote
	Prefs    []comsoc.Alternative `json:"prefs"`     //Préférences ordonnées de l'agent votant
	Options  []int                `json:"options"`   //Utilisé pour le seuil du vote par approbation
}

// Types utilisés pour la requête /result

type RequestResult struct {
	BallotId string `json:"ballot-id"` //id du scrutin dont on veut le résultat
}

type ResponseResult struct {
	//Objet renvoyé si code 200
	Winner  comsoc.Alternative   `json:"winner"`            //Alternative gagnante
	Ranking []comsoc.Alternative `json:"ranking,omitempty"` //Classement des alternatives (Champ facultatif)
}
