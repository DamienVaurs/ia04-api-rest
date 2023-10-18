package instances

import (
	"strconv"

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
func InitChap3Diap34(url string, n int, nbAlts int, listCin []chan []string, cout chan string) []restclientagent.RestClientAgent {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_" + strconv.Itoa(i+1)
	}

	res := make([]restclientagent.RestClientAgent, n)

	//On crée 3 types de scrutin pour les 3 méthodes de vote
	reqNewBallotMaj := restagent.RequestNewBallot{
		Rule:     restagent.Majority,
		Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{1, 2, 3, 4},
	}

	reqNewBallotBorda := restagent.RequestNewBallot{
		Rule:     restagent.Borda,
		Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{1, 2, 3, 4},
	}

	reqNewBallotCondorcet := restagent.RequestNewBallot{
		Rule:     restagent.Condorcet,
		Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{1, 2, 3, 4},
	}

	//5 agents votent pour [1, 2, 3, 4]
	for i := 0; i < 5; i++ {
		res[i] = *restclientagent.NewRestClientAgent(listAgentsId[i], url,
			reqNewBallotMaj,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{1, 2, 3, 4},
				Options:  nil,
			},
			listCin[i],
			cout)
	}

	//4 agents votent pour [1, 3, 2, 4]
	for i := 5; i < 9; i++ {
		res[i] = *restclientagent.NewRestClientAgent(listAgentsId[i], url,
			reqNewBallotMaj,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{1, 3, 2, 4},
				Options:  nil,
			},
			listCin[i],
			cout)
	}

	//2 agents votent pour [4, 2, 1, 3]
	for i := 9; i < 11; i++ {
		res[i] = *restclientagent.NewRestClientAgent(listAgentsId[i], url,
			reqNewBallotBorda,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{4, 2, 1, 3},
				Options:  nil,
			},
			listCin[i],
			cout)
	}

	//6 agents votent pour [4, 2, 3, 1]
	for i := 11; i < 17; i++ {
		res[i] = *restclientagent.NewRestClientAgent(listAgentsId[i], url,
			reqNewBallotBorda,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{4, 2, 3, 1},
				Options:  nil,
			},
			listCin[i],
			cout)
	}

	//8 agents votent pour [3, 2, 1, 4]
	for i := 17; i < 25; i++ {
		res[i] = *restclientagent.NewRestClientAgent(listAgentsId[i], url,
			reqNewBallotCondorcet,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{3, 2, 1, 4},
				Options:  nil,
			},
			listCin[i],
			cout)
	}

	//2 agents votent pour [4, 3, 2, 1]
	for i := 25; i < 27; i++ {
		res[i] = *restclientagent.NewRestClientAgent(listAgentsId[i], url,
			reqNewBallotCondorcet,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{4, 3, 2, 1},
				Options:  nil,
			},
			listCin[i],
			cout)
	}

	return res
}
