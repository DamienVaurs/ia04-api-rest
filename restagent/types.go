package restagent

import "gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"

type Request struct {
	Preferences []comsoc.Alternative `json:"pref"`
}

type Response struct {
	Result []comsoc.Alternative `json:"res"`
}

// Types utilisés pour la requête /new_ballot
type RequestNewBallot struct {
	Rule     string   `json:"rule"`
	Deadline string   `json:"deadline"`
	VoterIds []string `json:"voter-ids"`
	Alts     int      `json:"#alts"`
}

type ResponseNewBallot struct {
	//Objet renvoyé si code 201
	BallotId string `json:"ballot-id"`
}

// Types utilisés pour la requête /vote
type RequestVote struct {
	AgentId string               `json:"agent-id"`
	VoteId  string               `json:"vote-id"`
	Prefs   []comsoc.Alternative `json:"prefs"`
	Options []any                `json:"options"` //Voir si le type est ok
}

// Types utilisés pour la requête /result

type RequestResult struct {
	BallotId string `json:"ballot-id"`
}

type ResponseResult struct {
	//Objet renvoyé si code 200
	Winner  comsoc.Alternative   `json:"winner"`
	Ranking []comsoc.Alternative `json:"ranking"`
}
