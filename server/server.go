package server

import (
	"RexPromAgent/db"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"sync"
)

type Server struct {
	address    string
	router     *mux.Router
	db         *db.MySQLDB
	syncLocker sync.Mutex
}

func NewServer() Server {
	// Construct the sever router
	r := mux.NewRouter()
	// Construct database connection
	connectDB, err := db.ConnectDB()
	if err != nil {
		logrus.Fatalf("db connect error: %s", err)
	}
	// Construct http server
	s := Server{
		router:  r,
		db:      connectDB,
		address: viper.GetString("server.address"),
	}
	r.HandleFunc("/rules/alert/reload", s.SyncAlertRules).Methods("GET")
	return s
}

func (s Server) Start() {
	// Defer database close
	defer s.db.Close()
	// Start the http Server
	logrus.Infof("Server start at %v ....", s.address)
	logrus.Fatal(http.ListenAndServe(s.address, s.router))
}
