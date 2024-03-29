package restserveragent

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"
	"gitlab.utc.fr/milairhu/ia04-api-rest/restagent/endpoints"
)

type RestServerAgent struct {
	sync.Mutex                              //les requêtes doivent se faire l'une après l'autre, car certaines requêtes votent, d'autres demandent les résultats
	addr        string                      //adresse du serveur (ip:port)
	ballotsMap  map[string]comsoc.Profile   //associe l'id d'un ballot à son Profil
	ballotsList map[string]restagent.Ballot //associe l'id d'un ballot à son objet Ballot
	countBallot int                         //compteur de ballot (pout générer les ids)
}

func NewRestServerAgent(addr string) *RestServerAgent {
	b := make(map[string]comsoc.Profile, 0)
	l := make(map[string]restagent.Ballot, 0)
	return &RestServerAgent{addr: addr, ballotsMap: b, ballotsList: l, countBallot: 1}
}

// Test de la méthode (GET, POST, ...)
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc(endpoints.Results, rsa.doCalcResult)
	mux.HandleFunc(endpoints.Vote, rsa.doVote)
	mux.HandleFunc(endpoints.NewBallot, rsa.doCreateNewBallot)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
