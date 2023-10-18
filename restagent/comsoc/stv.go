package comsoc

/*
* Vote Simple Transférable (Single Transferable Vote (STV)
* Chaque individu indique donne son ordre de préférence
* Pour n candidats, on fait n−1 tours
* (à moins d’avoir avant une majorité stricte pour un candidat)
* On suppose qu’à chaque tour chaque individu “vote” pour son candidat
* préféré (parmi ceux encore en course)
* À chaque tour on élimine le plus mauvais candidat
* (celui qui a le moins de voix)
 */
//Le map fait +1 à chaque tour passé
func STV_SWF(p Profile) (Count, error) {
	ok := checkProfile(p)
	if ok != nil {
		return nil, ok
	}
	copyP := make(Profile, len(p)) //On copie le profil pour pouvoir faire des suppressions sans affecter l'original
	for i, votant := range p {
		copyP[i] = make([]Alternative, len(votant))
		copy(copyP[i], votant)
	}
	resMap := make(Count, len(p[0]))
	//on initialise le map à 0
	for _, alt := range copyP[0] {
		resMap[alt] = 0
	}

	for nbToursRestants := len(copyP[0]) - 1; nbToursRestants > 0; nbToursRestants-- {
		//On fait le tour on compte les voix
		//On élimine le plus mauvais candidat

		comptMap := make(Count, len(copyP[0]))
		for _, alt := range copyP[0] {
			comptMap[alt] = 0
		}
		for _, votant := range copyP {
			_, ok := comptMap[votant[0]]
			if ok {
				comptMap[votant[0]]++
			} else {
				comptMap[votant[0]] = 1
			}
		}
		//On a les scores de tous pour ce tour
		var miniCount int = len(copyP) + 1
		var miniAlt Alternative
		for alt, count := range comptMap {
			if count < miniCount {
				miniCount = count
				miniAlt = alt
			}
		}
		//On a le plus mauvais candidat, on le vire des votes
		for indP, votant := range copyP {
			for i, alt := range votant {
				if alt == miniAlt {
					votant[i] = votant[len(votant)-1]
				}
			}
			copyP[indP] = votant[:len(votant)-1]
		}
		//on incrémente chaque candidat passant au tour suivant
		for _, alt := range copyP[0] {
			if alt != miniAlt {
				resMap[alt]++
			}
		}

	}
	return resMap, nil
}

func STV_SCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := STV_SWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
