package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	address := flag.String("p", "0.0.0.0:8033", "set server port")
	help := flag.Bool("h", false, "this help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	hub := NewHub()
	sr := mux.NewRouter()

	sr.HandleFunc("/topic/", func(w http.ResponseWriter, r *http.Request) {
		createTopic(hub, w, r)
	})

	sr.HandleFunc("/topic/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		joinTopic(hub, w, r)
	})

	sr.PathPrefix("/").Handler(http.StripPrefix("", http.FileServer(http.Dir("./")))).Methods(http.MethodGet)

	log.Println("===============")
	log.Println("Listen Port", *address)
	log.Fatal(http.ListenAndServe(*address, sr))
}

