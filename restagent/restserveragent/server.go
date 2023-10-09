package restserveragent

import (
	"bytes"
	"encoding/json"
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
	sync.Mutex //les requêtes doivent se faire l'une après l'autre, car certaines requêtes votent, d'autres demandes les résultats
	id         string
	addr       string
	profile    comsoc.Profile
}

func NewRestServerAgent(addr string) *RestServerAgent {
	p := make(comsoc.Profile, 0)
	return &RestServerAgent{id: addr, addr: addr, profile: p}
}

// Test de la méthode
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (*RestServerAgent) decodeRequest(r *http.Request) (req restagent.Request, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *RestServerAgent) doCalc(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()
	// vérification de la méthode de la requête
	if !rsa.checkMethod("GET", w, r) {
		return
	}
	fmt.Println("Serveur recoit : ", r.URL)

	// traitement de la requête
	var resp restagent.Response

	// calcule de la réponse

	scf, err2 := comsoc.MajoritySCF(rsa.profile)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Error while processing result (SCF)"
		w.Write([]byte(msg))
		return
	}
	if len(scf) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Error : SCF ends with tie. No winner"
		w.Write([]byte(msg))
		return
	}
	resp.Result = scf //TODO : remplacer par un tie-breaker
	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (rsa *RestServerAgent) savePref(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()
	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Println("Serveur recoit : ", r.URL, req.Preferences)
	//Enregistre le vote
	rsa.profile = append(rsa.profile, req.Preferences)

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(req)
	w.Write(serial)
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc(endpoints.Results, rsa.doCalc)
	mux.HandleFunc(endpoints.Vote, rsa.savePref)

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
