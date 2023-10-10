package comsoc

import (
	"errors"
	"sort"
)

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt Alternative, prefs []Alternative) int {
	for i, a := range prefs {
		if a == alt {
			return i
		}
	}
	return -1
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	//return rank(alt1, prefs) < rank(alt2, prefs)
	for _, a := range prefs {
		if a == alt1 {
			return true
		}
		if a == alt2 {
			return false
		}
	}
	return false
}

// renvoie les meilleures alternatives pour un décompte donné
func maxCount(count Count) (bestAlts []Alternative) {
	//Si j'ai bien compris : retourner la ou les clés du map ayant les meilleurs valeurs
	//Etape 1 : parcourt pour cerner la meilleure valeur et enregistrement des meilleurs dans le tableau res
	var maxi = 0
	for i, v := range count {
		if v > maxi {
			maxi = v
			bestAlts = []Alternative{i}
		} else if v == maxi {
			bestAlts = append(bestAlts, i)
		}
	}
	return
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois par préférences
func checkProfile(prefs Profile) error {
	if len(prefs) < 1 {
		return errors.New("aucun vote n'a été soumis")
	}
	if len(prefs[0]) < 2 {
		return errors.New("moins de 2 candidats")
	}
	//Etape 1 : vérification de la complétude
	for _, v := range prefs {
		if len(v) != len(prefs[0]) {
			return errors.New("le profil n'est pas complet")
		}
	}
	//Etape 2 : vérification de l'unicité des alternatives
	for _, v := range prefs {
		for i, a := range v {
			for j, b := range v {
				if i != j && a == b {
					return errors.New("le profil n'est pas correct")
				}
			}
		}
	}
	return nil
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func checkProfileAlternative(prefs Profile, alts []Alternative) error {
	//Etape 1 :vérifie profil
	err := checkProfile(prefs)
	if err != nil {
		return err
	}

	//Etape 2 : vérifie que chaque alternative de alts apparaît exactement une fois par préférences
	for _, prof := range prefs {
		//Pour chaque profil
		for _, a := range alts {
			//Pour chaqe alternative de alts
			var isPresent = false
			for _, b := range prof {
				if a == b {
					isPresent = true
				}
			}
			if !isPresent {
				return errors.New("le profil n'est pas correct : il manque une alternative")
			}
		}
	}
	return nil
}

func MakeRanking(count Count) (ranking []Alternative) {

	//On parcourt le map
	ranking = make([]Alternative, len(count))
	//Comme les alternatives sont entre 1 et len(count) sans trou, on peut remplir naîvement la liste
	for i := 0; i < len(count); i++ {
		ranking[i] = Alternative(i + 1)
	}
	//On trie la liste en fonction du nombre de votes
	sort.Slice(ranking, func(i, j int) bool { return count[ranking[i]] > count[ranking[j]] })

	//TODO : vérifier que le classement est correct (ordre décroissant)
	return
}
