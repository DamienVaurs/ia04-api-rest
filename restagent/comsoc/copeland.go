package comsoc

import (
	"errors"
)

/*
* Règle de Copeland
* Le meilleur candidat est celui qui bat le plus d’autres candidats
* On associe à chaque candidat a le score suivant :
* pour chaque autre candidat b!= a
* +1 si une majorité préfère a à b,
* −1 si une majorité préfère b à a et
* 0 sinon
* Le candidat élu est celui qui a le plus haut score de Copeland
 */
func winCopelandDuel(p Profile, alt1 Alternative, alt2 Alternative) (int, error) {
	ok := checkProfile(p)
	if ok != nil {
		return -2, errors.New("profil non valide")
	}
	var i int = 0
	var nbWin = 0
	for _, votant := range p {
		for votant[i] != alt1 && votant[i] != alt2 {
			i++
		}
		if votant[i] == alt1 {
			nbWin++
		}
		i = 0
	}
	if nbWin > len(p)-nbWin {
		return 1, nil
	} else if nbWin < len(p)-nbWin {
		return -1, nil
	} else {
		return 0, nil
	}
}

func CopelandSWF(p Profile) (Count, error) {
	ok := checkProfile(p)
	if ok != nil {
		return nil, ok
	}
	var i int
	var j int
	resMap := make(Count, len(p[0]))
	for i = 0; i < len(p[0])-1; i++ {
		for j = i + 1; j < len(p[0]); j++ {
			win, _ := winCopelandDuel(p, p[0][i], p[0][j])
			_, ok1 := resMap[p[0][i]]
			if ok1 {
				resMap[p[0][i]] += win
			} else {
				resMap[p[0][i]] = win
			}

			_, ok2 := resMap[p[0][j]]
			if ok2 {
				resMap[p[0][j]] -= win
			} else {
				resMap[p[0][j]] = -win
			}
		}
	}
	return resMap, nil

}
func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := CopelandSWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
