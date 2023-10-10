package comsoc

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
		for i, alt := range votant {
			if i < thresholds[indVotant] {
				_, ok := count[alt]
				if !ok {
					count[alt] = 1
				} else {
					count[alt]++
				}
			}
		}
	}
	return count, nil
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, ranking []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)
	if err != nil {
		return nil, nil, err
	}
	return maxCount(count), MakeRanking(count), nil
}
