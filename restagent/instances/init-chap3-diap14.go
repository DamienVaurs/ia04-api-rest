package instances

import (
	"strconv"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

/**
* Diapo 14 du chapitre 3
* Méthode : Scrutin majoritaire simple
* Alternatives : 3
* 10 ordres de préférences : [1, 2, 3]
* 6 ordres de préférences : [2, 3, 1]
* 5 ordres de préférences : [3, 2, 1]
*
**/
func InitChap3Diap14(url string, n int, nbBallots int, nbAlts int, listCinVotants []chan []string, listCinBallots []chan []string, cout chan string) ([]restclientagent.RestClientVoteAgent, []restclientagent.RestClientBallotAgent) {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_vote_" + strconv.Itoa(i+1)
	}

	//Création du scrutin
	reqNewBallot := restagent.RequestNewBallot{
		Rule:     restagent.Majority,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{1, 2, 3},
	}

	voteAgents := make([]restclientagent.RestClientVoteAgent, n) // Liste des agents votants
	ballotAgent := []restclientagent.RestClientBallotAgent{      // Liste des agents tenannt un scrutin (un seul ici)
		*restclientagent.NewRestClientBallotAgent("ag_scrut", url, reqNewBallot, listCinBallots[0], cout),
	}

	//10 premiers agents votent pour [1, 2, 3]
	for i := 0; i < 10; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{1, 2, 3},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	//6 agents votent pour [2, 3, 1]
	for i := 10; i < 16; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{2, 3, 1},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	//5 agents votent pour [3, 2, 1]
	for i := 16; i < 21; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{3, 2, 1},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	return voteAgents, ballotAgent
}
