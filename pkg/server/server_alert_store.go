package server

import (
	"RexPromAgent/pkg/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (s Server) webhookPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read body error", err)
		return
	}

	data, err := ParsePromAlert(body)
	if err != nil {
		log.Printf("Invalid payload: %s\n", err)
		return
	}
	err = s.db.SaveAlert(data)
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

func ParsePromAlert(payload []byte) (*db.AlertGroup, error) {
	d := db.AlertGroup{}
	err := json.Unmarshal(payload, &d)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json webhook payload: %s", err)
	}
	return &d, nil
}
