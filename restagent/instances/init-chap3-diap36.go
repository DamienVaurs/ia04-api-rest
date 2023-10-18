package instances

import (
	"strconv"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

/**
* Diapo 36 du chapitre 3
* Méthode : Copeland
* Alternatives : 4
* 5 ordres de préférences : [1, 2, 3, 4]
* 4 ordres de préférences : [2, 3, 4, 1]
* 3 ordres de préférences : [4, 3, 1, 2]
*
**/
func InitChap3Diap36(url string, n int, nbBallots int, nbAlts int, listCinVotants []chan []string, listCinBallots []chan []string, cout chan string) ([]restclientagent.RestClientVoteAgent, []restclientagent.RestClientBallotAgent) {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_vote_" + strconv.Itoa(i+1)
	}

	//Création du scrutin
	reqNewBallot := restagent.RequestNewBallot{
		Rule:     restagent.Copeland,
		Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{1, 2, 3, 4},
	}

	voteAgents := make([]restclientagent.RestClientVoteAgent, n) // Liste des agents votants
	ballotAgent := []restclientagent.RestClientBallotAgent{      // Liste des agents tenannt un scrutin (un seul ici)
		*restclientagent.NewRestClientBallotAgent("ag_scrut", url, reqNewBallot, listCinBallots[0], cout),
	}

	//5 premiers agents votent pour [1, 2, 3, 4]
	for i := 0; i < 5; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{1, 2, 3, 4},
				Options:  nil,
			},
			listCinVotants[i],
			cout)
	}

	//4 agents votent pour [2, 3, 4, 1]
	for i := 5; i < 9; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{2, 3, 4, 1},
				Options:  nil,
			},
			listCinVotants[i],
			cout)
	}

	//3 agents votent pour [4, 3, 1, 2]
	for i := 9; i < 12; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{4, 3, 1, 2},
				Options:  nil,
			},
			listCinVotants[i],
			cout)
	}

	return voteAgents, ballotAgent
}
