package instances

import (
	"strconv"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

/**
* Teste de Approval avec et sans Tie-Break utilisé
* Méthode : Approval
* Alternatives : 5
* Scrutin 1 :
* 5 ordres de préférences : [1, 2, 3, 4, 5], seuil à 3
* 4 ordres de préférences : [2, 5, 3, 4, 1], seuil à 1
* 3 ordres de préférences : [3, 4, 1, 5, 2], seuil à 0
* 3 ordres de préférences : [3, 4, 1, 5, 2], seuil à 2
* Résultat attendu : [2, 3, 1, 4, 5]
*
* Scrutin 2 :
* 5 ordres de préférences : [1, 2, 3, 4, 5], seuil à 3
* 4 ordres de préférences : [2, 5, 3, 4, 1], seuil à 1
* 3 ordres de préférences : [3, 4, 1, 5, 2], seuil à 0
* Résultat attendu avec Tie-Break [5, 4, 3, 2, 1]: [2, 3, 1, 5, 4]
**/
func InitApproval(url string, n int, nbBallots int, nbAlts int, listCinVotants []chan []string, listCinBallots []chan []string, cout chan string) ([]restclientagent.RestClientVoteAgent, []restclientagent.RestClientBallotAgent) {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_vote_" + strconv.Itoa(i+1)
	}

	//Création des scrutins
	reqNewBallotNoTB := restagent.RequestNewBallot{
		Rule:     restagent.Approval,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{5, 4, 3, 2, 1},
	}
	reqNewBallotTB := restagent.RequestNewBallot{
		Rule:     restagent.Approval,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
		VoterIds: listAgentsId[:n-3], //On exlcu les 3 agents qui votent en plus au premier scrutin
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{5, 4, 3, 2, 1},
	}

	voteAgents := make([]restclientagent.RestClientVoteAgent, n) // Liste des agents votants
	ballotAgent := []restclientagent.RestClientBallotAgent{      // Liste des agents tenant un scrutin
		*restclientagent.NewRestClientBallotAgent("ag_scrut_No_Tie-Break", url, reqNewBallotNoTB, listCinBallots[0], cout),
		*restclientagent.NewRestClientBallotAgent("ag_scrut_Tie-Break", url, reqNewBallotTB, listCinBallots[1], cout),
	}

	//5 premiers agents votent pour [1, 2, 3, 4, 5] seuil à 3
	for i := 0; i < 5; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{1, 2, 3, 4, 5},
				Options: []int{3}, //seuil à 3
			},
			listCinVotants[i],
			cout)
	}

	//4 agents votent pour [2, 5, 3, 4, 1] seuil à 1
	for i := 5; i < 9; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{2, 5, 3, 4, 1},
				Options: []int{1}, //seuil à 1
			},
			listCinVotants[i],
			cout)
	}

	//3 agents votent pour [3, 4, 1, 5, 2] seuil à 0
	for i := 9; i < 12; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{3, 4, 1, 5, 2},
				Options: []int{0}, //seuil à 0
			},
			listCinVotants[i],
			cout)
	}

	//3 agents votent pour [3, 4, 1, 5, 2] seuil à 2
	for i := 12; i < 15; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{3, 4, 1, 5, 2},
				Options: []int{2}, //seuil à 2
			},
			listCinVotants[i],
			cout)
	}

	return voteAgents, ballotAgent
}
