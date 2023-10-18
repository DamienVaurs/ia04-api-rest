package instances

import (
	"math/rand"
	"strconv"

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

func Init10VotingAgents(url string, n int, nbAlts int, listCin []chan []string, cout chan string) []restclientagent.RestClientAgent {
	listAgentsId := make([]string, n)
	for i := 0; i < n; i++ {
		listAgentsId[i] = "ag_" + strconv.Itoa(i+1)
	}

	res := make([]restclientagent.RestClientAgent, n)

	//Agent 1
	res[0] = *restclientagent.NewRestClientAgent(listAgentsId[0], url,
		restagent.RequestNewBallot{
			Rule:     restagent.Majority,
			Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
			VoterIds: listAgentsId[:],
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
		},
		restagent.RequestVote{
			AgentId:  listAgentsId[0],
			BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
			Prefs:    generatePrefs(nbAlts),
			Options:  generateThresholds(nbAlts),
		},
		listCin[0],
		cout)

	//Agent 2
	res[1] = *restclientagent.NewRestClientAgent(listAgentsId[1], url,
		restagent.RequestNewBallot{
			Rule:     restagent.Borda,
			Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
			VoterIds: listAgentsId[:],
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
		},
		restagent.RequestVote{
			AgentId:  listAgentsId[1],
			BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
			Prefs:    generatePrefs(nbAlts),
			Options:  generateThresholds(nbAlts),
		},
		listCin[1],
		cout)

	//Agent 3
	res[2] = *restclientagent.NewRestClientAgent(listAgentsId[2], url,
		restagent.RequestNewBallot{
			Rule:     restagent.Approval,
			Deadline: "2024-12-31T23:59:59Z", //scrutin finissant dans trop longtemps pour avoir les résultats
			VoterIds: listAgentsId[:],
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
		},
		restagent.RequestVote{
			AgentId:  listAgentsId[2],
			BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
			Prefs:    generatePrefs(nbAlts),
			Options:  generateThresholds(nbAlts),
		},
		listCin[2],
		cout)

	// Agent 4
	res[3] = *restclientagent.NewRestClientAgent(listAgentsId[3], url,
		restagent.RequestNewBallot{
			Rule:     restagent.Majority,
			Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
			VoterIds: []string{},             //Pas de votants
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
		},
		restagent.RequestVote{
			AgentId:  listAgentsId[3],
			BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
			Prefs:    generatePrefs(nbAlts),
			Options:  generateThresholds(nbAlts),
		},
		listCin[3],
		cout)

	// Agent 5
	res[4] = *restclientagent.NewRestClientAgent(listAgentsId[4], url,
		restagent.RequestNewBallot{
			Rule:     restagent.Condorcet,
			Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
			VoterIds: listAgentsId[:],
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{5, 4, 3, 2, 1},
		},
		restagent.RequestVote{
			AgentId:  listAgentsId[4],
			BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
			Prefs:    generatePrefs(nbAlts),
			Options:  generateThresholds(nbAlts),
		},
		listCin[4],
		cout)

	// Agent 6
	res[5] = *restclientagent.NewRestClientAgent(listAgentsId[5], url,
		restagent.RequestNewBallot{
			Rule:     restagent.Copeland,
			Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
			VoterIds: listAgentsId[:],
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
		},
		restagent.RequestVote{
			AgentId:  listAgentsId[5],
			BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
			Prefs:    generatePrefs(nbAlts),
			Options:  generateThresholds(nbAlts),
		},
		listCin[5],
		cout)

	//  Agent 7
	res[6] = *restclientagent.NewRestClientAgent(listAgentsId[6], url,
		restagent.RequestNewBallot{
			Rule:     restagent.STV,
			Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
			VoterIds: listAgentsId[:],
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
		},
		restagent.RequestVote{
			AgentId:  listAgentsId[6],
			BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
			Prefs:    generatePrefs(nbAlts),
			Options:  generateThresholds(nbAlts),
		},
		listCin[6],
		cout)

	// Agent 8
	res[7] = *restclientagent.NewRestClientAgent(listAgentsId[7], url,
		restagent.RequestNewBallot{
			Rule:     restagent.Condorcet,
			Deadline: "2018-12-31T23:59:59Z",    //TODO : mettre une date cohérente quand on décommentera le code
			VoterIds: []string{listAgentsId[0]}, //1 seul votant, le résultat sera sa préférence
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
		},
		restagent.RequestVote{
			AgentId:  listAgentsId[7],
			BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
			Prefs:    generatePrefs(nbAlts),
			Options:  generateThresholds(nbAlts),
		},
		listCin[7],
		cout)

	// Agent 9
	res[8] = *restclientagent.NewRestClientAgent(listAgentsId[8], url,
		restagent.RequestNewBallot{
			Rule:     restagent.Approval,
			Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
			VoterIds: listAgentsId[:3],       //3 votants, le reste ne peut pas voter
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{5, 4, 3, 2, 1},
		},
		restagent.RequestVote{
			AgentId:  listAgentsId[8],
			BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
			Prefs:    generatePrefs(nbAlts),
			Options:  generateThresholds(nbAlts),
		},
		listCin[8],
		cout)

	// Agent 10
	res[9] = *restclientagent.NewRestClientAgent(listAgentsId[9], url,
		restagent.RequestNewBallot{
			Rule:     restagent.Borda,
			Deadline: "2018-12-31T23:59:59Z", //TODO : mettre une date cohérente quand on décommentera le code
			VoterIds: listAgentsId,
			Alts:     nbAlts,
			TieBreak: []comsoc.Alternative{1, 2, 3, 4, 5},
		},
		restagent.RequestVote{
			AgentId:  listAgentsId[9],
			BallotId: "", //Pas besoin de spécifier, l'agent vote pour le scrutin qu'il crée
			Prefs:    generatePrefs(nbAlts),
			Options:  generateThresholds(nbAlts),
		},
		listCin[9],
		cout)

	return res
}
