package instances

import (
	"strconv"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

/**
* Teste de Condorcet avec et sans gagnant
* Méthode : Condorcet
* Alternatives : 4
* Scrutin 1 :
* 5 ordres de préférences : [3, 1, 2, 4]
* 4 ordres de préférences : [2, 3, 4, 1]
* 3 ordres de préférences : [4, 3, 1, 2]
* Résultat attendu : 3 vainqueur de Condorcet
*
* Scrutin 2 :
* 5 ordres de préférences : [1, 2, 3, 4]
* 4 ordres de préférences : [2, 3, 4, 1]
* 3 ordres de préférences : [4, 3, 1, 2]
* Résultat attendu: pas de vainqueur de Condorcet
**/
func InitCondorcet(url string, n int, nbBallots int, nbAlts int, listCinVotants []chan []string, listCinBallots []chan []string, cout chan string) ([]restclientagent.RestClientVoteAgent, []restclientagent.RestClientBallotAgent) {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_vote_" + strconv.Itoa(i+1)
	}

	//Création des scrutins
	reqNewBallotNoWinner := restagent.RequestNewBallot{
		Rule:     restagent.Condorcet,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
		VoterIds: listAgentsId[5:],
		Alts:     nbAlts,
	}
	reqNewBallotWinner := restagent.RequestNewBallot{
		Rule:     restagent.Condorcet,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
		VoterIds: listAgentsId[:n-5],
		Alts:     nbAlts,
	}

	voteAgents := make([]restclientagent.RestClientVoteAgent, n) // Liste des agents votants
	ballotAgent := []restclientagent.RestClientBallotAgent{      // Liste des agents tenant un scrutin
		*restclientagent.NewRestClientBallotAgent("ag_scrut_No_Winner", url, reqNewBallotNoWinner, listCinBallots[0], cout),
		*restclientagent.NewRestClientBallotAgent("ag_scrut_Winner", url, reqNewBallotWinner, listCinBallots[1], cout),
	}

	//5 premiers agents votent pour [1, 2, 3, 4] pour le scrutin sans vainqueur
	for i := 0; i < 5; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{1, 2, 3, 4},
			},
			listCinVotants[i],
			cout)
	}

	//4 agents votent pour [2, 3, 4, 1]
	for i := 5; i < 9; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{2, 3, 4, 1},
			},
			listCinVotants[i],
			cout)
	}

	//3 agents votent pour [4, 3, 1, 2]
	for i := 9; i < 12; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{4, 3, 1, 2},
			},
			listCinVotants[i],
			cout)
	}
	//5 derniers agents votent pour [3, 1, 2, 4] pour le scrutin avec vainqueur
	for i := 12; i < 17; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{3, 1, 2, 4},
			},
			listCinVotants[i],
			cout)
	}

	return voteAgents, ballotAgent
}
