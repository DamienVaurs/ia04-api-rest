package restagent

//Ensemble des règles de vote prises en compte par l'application

const Approval = "approval"
const Borda = "borda"
const Condorcet = "condorcet"
const Copeland = "copeland"
const Majority = "majority"
const STV = "stv"

var Rules = []string{Approval, Borda, Condorcet, Copeland, Majority, STV}
