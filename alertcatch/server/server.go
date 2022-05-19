package server

import (
	"RexPromAgent/alertcatch/store"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	r     *mux.Router
	store *store.Storer
	debug bool
}

func NewServer(st *store.Storer, debug bool) Server {
	r := mux.NewRouter()

	s := Server{
		store: st,
		r:     r,
		debug: debug,
	}

	r.HandleFunc("/webhook", s.webhookPost).Methods("POST")

	return s
}

func (s Server) webhookPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read body error", err)
		return
	}

	data, err := store.ParsePromAlert(body)
	if err != nil {
		log.Printf("Invalid payload: %s\n", err)
		return
	}
	err = s.store.Conn.SaveAlert(data)
	if err != nil {
		log.Printf("failed to save alerts: %s\n", err)
		return
	}

	_, err = w.Write([]byte("saved"))
	if err != nil {
		log.Printf("error, response, %s", err)
		return
	}
}

func (s Server) Start(address string) {
	log.Println("Starting listener on", address, "using", s.store.Conn)
	log.Fatal(http.ListenAndServe(address, s.r))
}
