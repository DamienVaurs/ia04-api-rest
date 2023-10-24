package instances

import (
	"math/rand"
	"strconv"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

func generatePrefs(nbAlts int) []comsoc.Alternative {
	intPref := rand.Perm(nbAlts)
	altPref := make([]comsoc.Alternative, nbAlts)
	for i := 0; i < nbAlts; i++ {
		//conversion en []Alternative
		altPref[i] = comsoc.Alternative(intPref[i] + 1)
	}
	return altPref
}

func generateThresholds(nbAlts int) []int {
	//TODO : vérifier la valeur générée
	return []int{rand.Intn(5)}
}

func Init10VotingAgents(url string, n int, nbBallots int, nbAlts int, listCinVotants []chan []string, listCinBallots []chan []string, cout chan string) ([]restclientagent.RestClientVoteAgent, []restclientagent.RestClientBallotAgent) {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_vote_" + strconv.Itoa(i+1)
	}

	voteAgents := make([]restclientagent.RestClientVoteAgent, n)
	ballotAgents := make([]restclientagent.RestClientBallotAgent, nbBallots)

	//Création des scrutins
	for i, rule := range restagent.Rules {
		ballotAgents[i] = *restclientagent.NewRestClientBallotAgent("ag_scrut_"+rule, url,
			restagent.RequestNewBallot{
				Rule:     restagent.Rules[i],
				Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
				VoterIds: listAgentsId[:],
				Alts:     nbAlts,
				TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
			},
			listCinBallots[i],
			cout)
	}
	//Scrutin avec une deadline passée
	ballotAgents[nbBallots-2] = *restclientagent.NewRestClientBallotAgent("ag_scrut_passed_deadline", url,
		restagent.RequestNewBallot{
			Rule:     restagent.Majority,
			Deadline: "2018-12-31T23:59:59Z", //On met une deadline passée pour que les votants ne puisse pas voter pour ce scrutin
			VoterIds: listAgentsId[:],
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
		},
		listCinBallots[nbBallots-2],
		cout)

	//Scrutin sans votants autorisés
	ballotAgents[nbBallots-1] = *restclientagent.NewRestClientBallotAgent("ag_scrut_no_votant", url,
		restagent.RequestNewBallot{
			Rule:     restagent.Majority,
			Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
			VoterIds: []string{},
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
		},
		listCinBallots[nbBallots-1],
		cout)

	//Création des votants
	//Agent 1
	voteAgents[0] = *restclientagent.NewRestClientVoteAgent(listAgentsId[0], url,
		restagent.RequestVote{
			AgentId: listAgentsId[0],
			Prefs:   generatePrefs(nbAlts),
			Options: generateThresholds(nbAlts),
		},
		listCinVotants[0],
		cout)

	//Agent 2
	voteAgents[1] = *restclientagent.NewRestClientVoteAgent(
		listAgentsId[1], url,
		restagent.RequestVote{
			AgentId: listAgentsId[1],
			Prefs:   generatePrefs(nbAlts),
			Options: generateThresholds(nbAlts),
		}, listCinVotants[1], cout,
	)

	//Agent 3
	voteAgents[2] = *restclientagent.NewRestClientVoteAgent(listAgentsId[2], url,
		restagent.RequestVote{
			AgentId: listAgentsId[2],
			Prefs:   generatePrefs(nbAlts),
			Options: generateThresholds(nbAlts),
		},
		listCinVotants[2],
		cout)

	// Agent 4
	voteAgents[3] = *restclientagent.NewRestClientVoteAgent(listAgentsId[3], url,
		restagent.RequestVote{
			AgentId: listAgentsId[3],
			Prefs:   generatePrefs(nbAlts),
			Options: generateThresholds(nbAlts),
		},
		listCinVotants[3],
		cout)

	// Agent 5
	voteAgents[4] = *restclientagent.NewRestClientVoteAgent(listAgentsId[4], url,
		restagent.RequestVote{
			AgentId: listAgentsId[4],
			Prefs:   generatePrefs(nbAlts),
			Options: generateThresholds(nbAlts),
		},
		listCinVotants[4],
		cout)

	// Agent 6
	voteAgents[5] = *restclientagent.NewRestClientVoteAgent(listAgentsId[5], url,
		restagent.RequestVote{
			AgentId: listAgentsId[5],
			Prefs:   generatePrefs(nbAlts),
			Options: generateThresholds(nbAlts),
		},
		listCinVotants[5],
		cout)

	//  Agent 7
	voteAgents[6] = *restclientagent.NewRestClientVoteAgent(listAgentsId[6], url,
		restagent.RequestVote{
			AgentId: listAgentsId[6],
			Prefs:   generatePrefs(nbAlts),
			Options: generateThresholds(nbAlts),
		},
		listCinVotants[6],
		cout)

	// Agent 8
	voteAgents[7] = *restclientagent.NewRestClientVoteAgent(listAgentsId[7], url,
		restagent.RequestVote{
			AgentId: listAgentsId[7],
			Prefs:   generatePrefs(nbAlts),
			Options: generateThresholds(nbAlts),
		},
		listCinVotants[7],
		cout)

	// Agent 9
	voteAgents[8] = *restclientagent.NewRestClientVoteAgent(listAgentsId[8], url,
		restagent.RequestVote{
			AgentId: listAgentsId[8],
			Prefs:   generatePrefs(nbAlts),
			Options: generateThresholds(nbAlts),
		},
		listCinVotants[8],
		cout)

	// Agent 10
	voteAgents[9] = *restclientagent.NewRestClientVoteAgent(listAgentsId[9], url,
		restagent.RequestVote{
			AgentId: listAgentsId[9],
			Prefs:   generatePrefs(nbAlts),
			Options: generateThresholds(nbAlts),
		},
		listCinVotants[9],
		cout)

	return voteAgents, ballotAgents
}
