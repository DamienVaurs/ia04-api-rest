package comsoc

import "errors"

//Est-ce qu'on gagne si on est égalité?
func winDuel(p Profile, alt1 Alternative, alt2 Alternative) (bool, error) {
	ok := checkProfile(p)
	if ok != nil {
		return false, errors.New("profil non valide")
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
	return nbWin > len(p)-nbWin, nil //Voir si inégalité stricte ou pas
}

//Donne le gagnant de Codorcet ou nil si il n'y en a pas
func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	ok := checkProfile(p)
	if ok != nil {
		return nil, errors.New("Profile non valide")
	}
	var m int = len(p[0]) //nb alternatives
	var n int = len(p)    //nb individus

	//Cas particuliers
	if m == 1 {
		//Si une seule alternative
		bestAlts = []Alternative{p[0][0]}
		return bestAlts, nil
	}
	if n == 1 {
		//Si un seul individu
		return []Alternative{p[0][0]}, nil
	}

	//Cas général
	//On fait tous les duels. On voit si un gagne tous ses duels.
	//Si oui, c'est le gagnant de Condorcet
	//Si non, il n'y a pas de gagnant de Condorcet
	var i int = 0
	var j int = 0
	for i = 0; i < m; i++ {
		var nbWin int = 0
		for j = 0; j < m; j++ {
			if i != j {
				win, _ := winDuel(p, p[0][i], p[0][j])
				if win {
					nbWin++
				}
			}
		}
		if nbWin == m-1 {
			return []Alternative{p[0][i]}, nil
		}
	}
	return []Alternative{}, nil
}
