package comsoc

import "fmt"

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}
	count = make(Count, len(p[0])) //initialisation du map
	for _, alt := range p[0] {
		//On initialise à 0
		count[alt] = 0
	}
	//Recensement des votes de tous les profils, de 0 à tresholds[i]
	for indVotant, votant := range p {
		//Pour tout les votes recencés
		for i, alt := range votant {
			//Pour chaque alternative d'un vote
			if thresholds[indVotant] < 0 || thresholds[indVotant] > len(votant) {
				return nil, fmt.Errorf("threshold %d is incorrect with %d alternatives", thresholds[indVotant], len(votant))
			}
			if i < thresholds[indVotant] {
				//Si on compte cette alternative
				_, ok := count[alt]
				if !ok {
					count[alt] = 1
				} else {
					count[alt]++
				}
			} else {
				break
			}
		}
	}
	return count, nil
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
