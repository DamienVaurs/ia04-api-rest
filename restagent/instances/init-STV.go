package instances

import (
	"strconv"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

/**
* Teste de STV avec et sans Tie-Break utilisé
* Méthode : STV
* Alternatives : 4
* Scrutin 1 :
* 5 ordres de préférences : [1, 2, 3, 4]
* 4 ordres de préférences : [2, 3, 4, 1]
* 3 ordres de préférences : [4, 3, 1, 2]
* Résultat attendu : [1, 2, 4, 3]
*
* Scrutin 2 :
* 3 ordres de préférences : [1, 2, 3, 4]
* 3 ordres de préférences : [2, 3, 4, 1]
* 3 ordres de préférences : [4, 3, 1, 2]
* Résultat attendu avec Tie-Break [4, 3, 2, 1]: [2, 4, 1, 3]
**/
func InitSTV(url string, n int, nbBallots int, nbAlts int, listCinVotants []chan []string, listCinBallots []chan []string, cout chan string) ([]restclientagent.RestClientVoteAgent, []restclientagent.RestClientBallotAgent) {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_vote_" + strconv.Itoa(i+1)
	}

	//Création des scrutins
	reqNewBallotNoTB := restagent.RequestNewBallot{
		Rule:     restagent.STV,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{4, 3, 2, 1},
	}
	reqNewBallotTB := restagent.RequestNewBallot{
		Rule:     restagent.STV,
		Deadline: time.Now().Add(5 * time.Second).Format(time.RFC3339),
		VoterIds: listAgentsId[:n-3], //On exlcu les 3 agents qui votent en plus au premier scrutin
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{4, 3, 2, 1},
	}

	voteAgents := make([]restclientagent.RestClientVoteAgent, n) // Liste des agents votants
	ballotAgent := []restclientagent.RestClientBallotAgent{      // Liste des agents tenant un scrutin
		*restclientagent.NewRestClientBallotAgent("ag_scrut_No_Tie-Break", url, reqNewBallotNoTB, listCinBallots[0], cout),
		*restclientagent.NewRestClientBallotAgent("ag_scrut_Tie-Break", url, reqNewBallotTB, listCinBallots[1], cout),
	}

	//3 premiers agents votent pour [1, 2, 3, 4]
	for i := 0; i < 3; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{1, 2, 3, 4},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}
	//2 autres agents ne votants qu'au scrutin 1 ont les mêmes choix :
	for i := n - 3; i < n-1; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{1, 2, 3, 4},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	//3 agents votent pour [2, 3, 4, 1]
	for i := 3; i < 6; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{2, 3, 4, 1},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}
	//1 autre agent ne votant qu'au scrutin 1 a les mêmes choix :

	voteAgents[n-1] = *restclientagent.NewRestClientVoteAgent(listAgentsId[n-1], url,
		restagent.RequestVote{
			AgentId: listAgentsId[n-1],
			Prefs:   []comsoc.Alternative{2, 3, 4, 1},
			Options: nil,
		},
		listCinVotants[n-1],
		cout)

	//3 agents votent pour [4, 3, 1, 2]
	for i := 6; i < 9; i++ {
		voteAgents[i] = *restclientagent.NewRestClientVoteAgent(listAgentsId[i], url,
			restagent.RequestVote{
				AgentId: listAgentsId[i],
				Prefs:   []comsoc.Alternative{4, 3, 1, 2},
				Options: nil,
			},
			listCinVotants[i],
			cout)
	}

	return voteAgents, ballotAgent
}
