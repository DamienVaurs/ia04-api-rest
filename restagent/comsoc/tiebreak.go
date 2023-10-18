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

// Remarque : obligé de créer une fonction SWF avec Tie-break particulière pour approval car il faut prendre en compte le seuil
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

// Remarque : obligé de créer une fonction SWF avec Tie-break particulière pour STV car le départage est différent. On utilise le Tie-Break au sein même de l'algorithme
// TODO : vérifier que ça marche
func STV_SWF_TieBreak(p Profile, tieBreak []Alternative) (Count, error) {
	ok := checkProfile(p)
	if ok != nil {
		return nil, ok
	}
	copyP := make(Profile, len(p)) //On copie le profil pour pouvoir faire des suppressions sans affecter l'original
	for i, votant := range p {
		copyP[i] = make([]Alternative, len(votant))
		copy(copyP[i], votant)
	}

	tieBreakMap := make(map[Alternative]int, len(tieBreak))
	for i, alt := range tieBreak {
		tieBreakMap[alt] = i
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
		miniAlts := make([]Alternative, 0)
		for alt, count := range comptMap {
			if count < miniCount {
				miniCount = count
				miniAlts = append(miniAlts, alt)
			}
		}
		//On a les plus mauvais candidats, on en vire un des votes selon le Tie-break fourni
		var miniAlt Alternative
		miniValInTieBreak := len(tieBreak) + 1
		for _, alt := range miniAlts {
			if tieBreakMap[alt] < miniValInTieBreak {
				miniValInTieBreak = tieBreakMap[alt]
				miniAlt = alt
			}
		}
		//On a désigné le candidat dont il faut se débarasser
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
