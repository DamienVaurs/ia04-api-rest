package comsoc

import (
	"errors"
)

///// Tie breakers
/**
 * En cas d'égalité, les SCF doivent renvoyer un seul élement
 * et les SWF doivent renvoyer un ordre total sans égalité.
 * On utilisera pour cela des fonctions de tie-break qui
 * étant donné un ensemble d'alternatives, renvoie la meilleure.
 * Elles respectent la signature suivante
 * (une erreur pouvant se produire si le slice d'alternatives est
 * vide) :
**/
func TieBreakFactory(ordreStrict []Alternative) func([]Alternative) (Alternative, error) {
	//Les alternatives fournies servent à départager les ex aequo -> ordre strict
	return func(alts []Alternative) (Alternative, error) {
		//On considère que les alternatives fournies sont de bases toutes à égalité
		if len(alts) == 0 {
			return -1, errors.New("aucune alternative fournie")
		} else {
			order := make(map[Alternative]int, len(ordreStrict))
			for i, alt := range ordreStrict {
				order[alt] = len(ordreStrict) - i
			}
			//le map order contient les alternatives associées à leur rang
			var maxVal int
			var maxAlt Alternative
			for _, alt := range alts {
				if order[alt] > maxVal {
					maxVal = order[alt]
					maxAlt = alt
				}
			}
			return maxAlt, nil
		}
	}
}

// Pour avoir des SWF sans ex aequo
func SWFFactory(swf func(p Profile) (Count, error), tieBreaker func([]Alternative) (Alternative, error)) func(Profile) ([]Alternative, error) {
	//Retourne une fonction qui retourne les alternatives ordonnées
	return func(p Profile) ([]Alternative, error) {
		count, err := swf(p)
		if err != nil {
			return nil, err
		}
		res := make([]Alternative, len(count))
		// On remplit res avec les alternatives : plus count[alt] est grand, plus alt est bien classée

		//Idée du prof : on multiplie toutes les alternatives par
		//nbAlt, et pour chaque alternative ex aequo,
		//on fait +1, +2 etc pour départager
		// -> ici, on fait pas ça
		invCount := make(map[int][]Alternative, len(count)) //dico {score : [candidats]}
		var maxScore int
		var minScore int
		for alt, score := range count {
			//On remplit le dictionnaire invCount et on enregistre les scores max et min
			invCount[score] = append(invCount[score], alt)
			if score > maxScore {
				maxScore = score
			} else if score < minScore {
				minScore = score
			}

		}
		var currIndex = 0
		for i := maxScore; i >= minScore; i-- {
			tab, ok := invCount[i]
			if ok {
				//Si on a des candidats correspondant à ce score,
				//on les trie en fonction du tiebreak et on les ajoute
				//au tableau res
				for len(tab) > 1 {
					//tant qu'il y a plusieurs éléments égalité,
					//on retire le meilleur de la liste
					//et on l'ajoute à res
					best, err := tieBreaker(tab)
					if err != nil {
						return nil, err
					}
					res[currIndex] = best
					currIndex++
					//On supprime l'élément de tab
					for i, alt := range tab {
						if alt == best {
							tab[i] = tab[len(tab)-1]
							tab = tab[:len(tab)-1]
							break
						}
					}
				}
				if len(tab) == 1 {
					//On ajoute le dernier élément au tableau
					res[currIndex] = tab[0]
					currIndex++
				}
			}
		}
		return res, nil
	}
}
func SCFFactory(scf func(p Profile) ([]Alternative, error), tieBreaker func([]Alternative) (Alternative, error)) func(Profile) (Alternative, error) {
	//Applique la fonction scf sur le profile puis départage les ex aequo avec tieBreaker. Renvoie la fonction qui applique scf mais sans ex aequo
	return func(p Profile) (Alternative, error) {
		bestAlts, err := scf(p)
		if err != nil {
			return -1, err
		}
		//On a les meilleures alternatives. On utilise tiebreaker pour départager
		return tieBreaker(bestAlts)
	}
}

func MakeApprovalRankingWithTieBreak(p Profile, threshold []int, tieBreaker func([]Alternative) (Alternative, error)) ([]Alternative, error) {
	count, err := ApprovalSWF(p, threshold)
	if err != nil {
		return nil, err
	}
	res := make([]Alternative, len(count))
	// On remplit res avec les alternatives : plus count[alt] est grand, plus alt est bien classée

	//Idée du prof : on multiplie toutes les alternatives par
	//nbAlt, et pour chaque alternative ex aequo,
	//on fait +1, +2 etc pour départager
	// -> ici, on fait pas ça
	invCount := make(map[int][]Alternative, len(count)) //dico {score : [candidats]}
	var maxScore int
	var minScore int
	for alt, score := range count {
		//On remplit le dictionnaire invCount et on enregistre les scores max et min
		invCount[score] = append(invCount[score], alt)
		if score > maxScore {
			maxScore = score
		} else if score < minScore {
			minScore = score
		}

	}
	var currIndex = 0
	for i := maxScore; i >= minScore; i-- {
		tab, ok := invCount[i]
		if ok {
			//Si on a des candidats correspondant à ce score,
			//on les trie en fonction du tiebreak et on les ajoute
			//au tableau res
			for len(tab) > 1 {
				//tant qu'il y a plusieurs éléments égalité,
				//on retire le meilleur de la liste
				//et on l'ajoute à res
				best, err := tieBreaker(tab)
				if err != nil {
					return nil, err
				}
				res[currIndex] = best
				currIndex++
				//On supprime l'élément de tab
				for i, alt := range tab {
					if alt == best {
						tab[i] = tab[len(tab)-1]
						tab = tab[:len(tab)-1]
						break
					}
				}
			}
			if len(tab) == 1 {
				//On ajoute le dernier élément au tableau
				res[currIndex] = tab[0]
				currIndex++
			}
		}
	}
	return res, nil

}
