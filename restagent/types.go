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
	Rule     string               `json:"rule"`
	Deadline string               `json:"deadline"`
	VoterIds []string             `json:"voter-ids"`
	Alts     int                  `json:"#alts"`
	TieBreak []comsoc.Alternative `json:"tie-break"`
}

type ResponseNewBallot struct {
	//Objet renvoyé si code 201
	BallotId string `json:"ballot-id"`
}

// Type utilisé pour la requête /vote
type RequestVote struct {
	AgentId  string               `json:"agent-id"`
	BallotId string               `json:"ballot-id"`
	Prefs    []comsoc.Alternative `json:"prefs"`
	Options  []int                `json:"options"` //Utilisé pour le seuil du vote par approbation
}

// Types utilisés pour la requête /result

type RequestResult struct {
	BallotId string `json:"ballot-id"`
}

type ResponseResult struct {
	//Objet renvoyé si code 200
	Winner  comsoc.Alternative   `json:"winner"`
	Ranking []comsoc.Alternative `json:"ranking,omitempty"` //Champ facultatif
}
