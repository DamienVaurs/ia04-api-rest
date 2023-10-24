package instances

import (
	"strconv"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

/**
* Diapo 34 du chapitre 3
* Méthodes : Scrutin majoritaire simple,
*			 Scrutin majoritaire à deux tours (non implémenté),
*			 Borda
*            Condorcet
* Alternatives : 4
* 5 ordres de préférences : [1, 2, 3, 4]
* 4 ordres de préférences : [1, 3, 2, 4]
* 2 ordres de préférences : [4, 2, 1, 3]
* 6 ordres de préférences : [4, 2, 3, 1]
* 8 ordres de préférences : [3, 2, 1, 4]
* 2 ordres de préférences : [4, 3, 2, 1]
*
**/
func InitChap3Diap34(url string, n int, nbBallot, nbAlts int, listCinVotants []chan []string, listCinBallots []chan []string, cout chan string) ([]restclientagent.RestClientVoteAgent, []restclientagent.RestClientBallotAgent) {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_vote_" + strconv.Itoa(i+1)
	}

	//On crée 3 scrutins, pour les 3 méthodes de vote
	reqNewBallotMaj := restagent.RequestNewBallot{
		Rule:     restagent.Majority,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339), //on laisse 5s pour que les votes se clôturent
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{1, 2, 3, 4},
	}

	reqNewBallotBorda := restagent.RequestNewBallot{
		Rule:     restagent.Borda,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{1, 2, 3, 4},
	}

	reqNewBallotCondorcet := restagent.RequestNewBallot{
		Rule:     restagent.Condorcet,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{1, 2, 3, 4},
	}

	voteAgents := make([]restclientagent.RestClientVoteAgent, n)     // Liste des agents votants
	ballotAgents := make([]restclientagent.RestClientBallotAgent, 3) // Liste des agents tenant un scrutin (un seul ici)

	//Création des scrutins
	ballotAgents[0] = *restclientagent.NewRestClientBallotAgent("ag_scrut_maj", url, reqNewBallotMaj, listCinBallots[0], cout)
	ballotAgents[1] = *restclientagent.NewRestClientBallotAgent("ag_scrut_borda", url, reqNewBallotBorda, listCinBallots[1], cout)
	ballotAgents[2] = *restclientagent.NewRestClientBallotAgent("ag_scrut_condorcet", url, reqNewBallotCondorcet, listCinBallots[2], cout)

	//5 agents votent pour [1, 2, 3, 4]
	for i := 0; i < 5; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{1, 2, 3, 4},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	//4 agents votent pour [1, 3, 2, 4]
	for i := 5; i < 9; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{1, 3, 2, 4},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	//2 agents votent pour [4, 2, 1, 3]
	for i := 9; i < 11; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{4, 2, 1, 3},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	//6 agents votent pour [4, 2, 3, 1]
	for i := 11; i < 17; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{4, 2, 3, 1},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	//8 agents votent pour [3, 2, 1, 4]
	for i := 17; i < 25; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{3, 2, 1, 4},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	//2 agents votent pour [4, 3, 2, 1]
	for i := 25; i < 27; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{4, 3, 2, 1},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	return voteAgents, ballotAgents
}
