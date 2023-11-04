package instances

import (
	"math/rand"
	"strconv"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

func InitVotingAgents(url string, n int, nbBallots int, nbAlts int, listCinVotants []chan []string, listCinBallots []chan []string, cout chan string) ([]restclientagent.RestClientVoteAgent, []restclientagent.RestClientBallotAgent) {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_vote_" + strconv.Itoa(i+1)
	}

	voteAgents := make([]restclientagent.RestClientVoteAgent, n)
	ballotAgents := make([]restclientagent.RestClientBallotAgent, nbBallots)

	//Création des scrutins

	for i := 0; i < nbBallots; i++ {
		ballotAgents[i] = *restclientagent.NewRestClientBallotAgent("ag_scrut_"+strconv.Itoa(i+1), url,
			restagent.RequestNewBallot{
				Rule:     restagent.Rules[rand.Intn(len(restagent.Rules))],
				Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
				VoterIds: listAgentsId[:],
				Alts:     nbAlts,
				TieBreak: generatePrefs(nbAlts),
			},
			listCinBallots[i],
			cout)
	}

	//Création des votants

	for i := 0; i < n; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   generatePrefs(nbAlts),
				Options: generateThresholds(nbAlts),
			},
			listCinVotants[i],
			cout)
	}

	return voteAgents, ballotAgents
}
