package instances

import (
	"strconv"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

/**
* Diapo 12 et 13 du chapitre 3
* Méthode : Trouver si gagnant de Condorcet
* Alternatives : 3
* Profil 1 :[1, 2, 3], [1, 3, 2], [3, 2, 1]
* Profil 2 :[1, 2, 3], [2, 3, 1], [3, 1, 2]
*
**/
func InitChap3Diap12_13(url string, n int, nbBallots int, nbAlts int, listCinVotants []chan []string, listCinBallots []chan []string, cout chan string) ([]restclientagent.RestClientVoteAgent, []restclientagent.RestClientBallotAgent) {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_vote_" + strconv.Itoa(i+1)
	}

	//Création des scrutins
	reqNewBallot12 := restagent.RequestNewBallot{
		Rule:     restagent.Condorcet,
		Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
		VoterIds: listAgentsId[:3],
		Alts:     nbAlts,
	}
	reqNewBallot13 := restagent.RequestNewBallot{
		Rule:     restagent.Condorcet,
		Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
		VoterIds: listAgentsId[3:n],
		Alts:     nbAlts,
	}

	voteAgents := make([]restclientagent.RestClientVoteAgent, n) // Liste des agents votants
	ballotAgent := []restclientagent.RestClientBallotAgent{      // Liste des agents tenant un scrutin (deux ici)
		*restclientagent.NewRestClientBallotAgent("ag_scrut_CondExist", url, reqNewBallot12, listCinBallots[0], cout),
		*restclientagent.NewRestClientBallotAgent("ag_scrut_NotCond", url, reqNewBallot13, listCinBallots[1], cout),
	}

	//Votants pour premier scrutin
	voteAgents[0] = *restclientagent.NewRestClientVoteAgent(listAgentsId[0], url,
		restagent.RequestVote{
			AgentId:  listAgentsId[0],
			BallotId: "ag_scrut_CondExist",
			Prefs:    []comsoc.Alternative{1, 2, 3},
			Options:  nil,
		},
		listCinVotants[0],
		cout)

	voteAgents[1] = *restclientagent.NewRestClientVoteAgent(listAgentsId[1], url,
		restagent.RequestVote{
			AgentId:  listAgentsId[1],
			BallotId: "ag_scrut_CondExist",
			Prefs:    []comsoc.Alternative{1, 3, 2},
			Options:  nil,
		},
		listCinVotants[1],
		cout)
	voteAgents[2] = *restclientagent.NewRestClientVoteAgent(listAgentsId[2], url,
		restagent.RequestVote{
			AgentId:  listAgentsId[2],
			BallotId: "ag_scrut_CondExist",
			Prefs:    []comsoc.Alternative{3, 2, 1},
			Options:  nil,
		},
		listCinVotants[2],
		cout)

	//Votants pour second scrutin
	voteAgents[3] = *restclientagent.NewRestClientVoteAgent(listAgentsId[3], url,
		restagent.RequestVote{
			AgentId:  listAgentsId[3],
			BallotId: "ag_scrut_NotCond",
			Prefs:    []comsoc.Alternative{1, 2, 3},
			Options:  nil,
		},
		listCinVotants[3],
		cout)

	voteAgents[4] = *restclientagent.NewRestClientVoteAgent(listAgentsId[4], url,
		restagent.RequestVote{
			AgentId:  listAgentsId[4],
			BallotId: "ag_scrut_NotCond",
			Prefs:    []comsoc.Alternative{2, 3, 1},
			Options:  nil,
		},
		listCinVotants[4],
		cout)

	voteAgents[5] = *restclientagent.NewRestClientVoteAgent(listAgentsId[5], url,
		restagent.RequestVote{
			AgentId:  listAgentsId[5],
			BallotId: "ag_scrut_NotCond",
			Prefs:    []comsoc.Alternative{3, 1, 2},
			Options:  nil,
		},
		listCinVotants[5],
		cout)

	return voteAgents, ballotAgent
}