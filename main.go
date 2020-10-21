//go:generate go run data/generate.go

package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"filegogo/data"
	"filegogo/lightcable"

	"github.com/gorilla/mux"
)

func main() {

	address := flag.String("p", "0.0.0.0:8033", "set server port")
	configPath := flag.String("c", "./config.json", "use config.json")
	help := flag.Bool("h", false, "this help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	hub := lightcable.NewHub()
	sr := mux.NewRouter()

	sr.HandleFunc("/topic/", func(w http.ResponseWriter, r *http.Request) {
		lightcable.CreateTopic(hub, w, r)
	})

	sr.HandleFunc("/topic/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		lightcable.JoinTopic(hub, w, r)
	})

	sr.HandleFunc("/config.json", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Read config: %s", *configPath)

		w.Header().Add("Content-type", "application/json")
		file, err := os.Open(*configPath)
		if err != nil {
			return
		}
		_, err = io.Copy(w, file)
		if err != nil {
			return
		}
	})

	sr.PathPrefix("/").Handler(http.StripPrefix("", http.FileServer(data.Dir))).Methods(http.MethodGet)

	log.Println("===============")
	log.Println("Listen Port", *address)
	log.Fatal(http.ListenAndServe(*address, sr))
}
