package comsoc

//Méthode Borda
func BordaSWF(p Profile) (Count, error) {
	err := checkProfile(p)
	if err != nil {
		return nil, err
	}
	count := make(Count, len(p[0])) //initialisation du map
	for _, alt := range p[0] {
		count[alt] = 0
	}
	//Recensement des votes du profil
	var nbAlt = len(p[0])
	for _, votant := range p {
		for i, alt := range votant {
			_, ok := count[alt]
			if !ok {
				//En réalité, pour le premier votant, on passe toujours dan ce if
				count[alt] = nbAlt - 1 - i
			} else {
				//En réalité, SAUF pour le premier votant, on passe toujours dan ce if
				count[alt] += nbAlt - 1 - i
			}
		}
	}
	return count, nil

}

func BordaSCF(p Profile) (bestAlts []Alternative, ranking []Alternative, err error) {
	count, err := BordaSWF(p)
	if err != nil {
		return nil, nil, err
	}
	return maxCount(count), MakeRanking(count), nil
}
