package comsoc

//Méthode majorité simple

func MajoritySWF(p Profile) (count Count, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}
	count = make(Count, len(p[0])) //initialisation du map
	for _, alt := range p[0] {
		//On initialise à 0
		count[alt] = 0
	}
	//Recensement des votes du profil
	for _, votant := range p {
		_, ok := count[votant[0]] //votant[0] est le préféré de votant
		if ok {
			count[votant[0]]++
		} else {
			count[votant[0]] = 1
		}
	}
	return count, nil

}

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := MajoritySWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
