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

	hub := newHub()
	r := mux.NewRouter()

	r.HandleFunc("/ws/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Println("===============")
	log.Println("Listen Port", *address)
	log.Fatal(http.ListenAndServe(*address, r))
}

