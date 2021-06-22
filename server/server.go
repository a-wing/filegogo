package server

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"filegogo/lightcable"

	"github.com/gorilla/mux"
)

//go:embed dist
var dist embed.FS

func Run(address, configPath string) {
	cable := lightcable.NewServer()
	sr := mux.NewRouter()

	sr.HandleFunc("/"+lightcable.PrefixShare+"/", func(w http.ResponseWriter, r *http.Request) {
		cable.JoinTopic(w, r)
	})

	sr.HandleFunc("/"+lightcable.PrefixShare+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		cable.JoinTopic(w, r)
	})

	sr.HandleFunc("/config.json", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Read config: %s", configPath)

		w.Header().Add("Content-type", "application/json")
		file, err := os.Open(configPath)
		if err != nil {
			return
		}
		_, err = io.Copy(w, file)
		if err != nil {
			return
		}
	})

	fsys, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatal(err)
	}
	sr.PathPrefix("/").Handler(http.StripPrefix("", http.FileServer(http.FS(fsys)))).Methods(http.MethodGet)

	log.Println("===============")
	log.Println("Listen Port", address)
	log.Fatal(http.ListenAndServe(address, sr))
}
