package instances

import (
	"strconv"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/restclientagent"
)

/**
* Diapo 14 du chapitre 3
* Méthode : Scrutin majoritaire simple
* Alternatives : 3
* 10 ordres de préférences : [1, 2, 3]
* 6 ordres de préférences : [2,3,1]
* 5 ordres de préférences : [3,2,1]
*
**/
func InitChap3Diap14(url string, n int, nbAlts int, listCin []chan []string, cout chan string) []restclientagent.RestClientAgent {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_" + strconv.Itoa(i+1)
	}

	res := make([]restclientagent.RestClientAgent, n)

	//Tout le monde crée des scrutins similaires
	reqNewBallot := restagent.RequestNewBallot{
		Rule:     restagent.Majority,
		Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
		VoterIds: listAgentsId[:],
		Alts:     nbAlts,
		TieBreak: []comsoc.Alternative{1, 2, 3},
	}

	//10 premiers agents votent pour [1, 2, 3]
	for i := 0; i < 10; i++ {
		res[i] = *restclientagent.NewRestClientAgent(listAgentsId[i], url,
			reqNewBallot,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{1, 2, 3},
				Options:  nil,
			},
			listCin[i],
			cout)
	}

	//6 agents votent pour [2, 3, 1]
	for i := 10; i < 16; i++ {
		res[i] = *restclientagent.NewRestClientAgent(listAgentsId[i], url,
			reqNewBallot,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{2, 3, 1},
				Options:  nil,
			},
			listCin[i],
			cout)
	}

	//5 agents votent pour [3, 2, 1]
	for i := 16; i < 21; i++ {
		res[i] = *restclientagent.NewRestClientAgent(listAgentsId[i], url,
			reqNewBallot,
			restagent.RequestVote{
				AgentId:  listAgentsId[i],
				BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
				Prefs:    []comsoc.Alternative{3, 2, 1},
				Options:  nil,
			},
			listCin[i],
			cout)
	}

	return res
}
