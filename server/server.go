package server

import (
	"context"
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
	sr := mux.NewRouter()

	cable := lightcable.NewServer()
	go cable.Run(context.Background())
	httpServer := NewServer(cable)

	sr.HandleFunc("/"+PrefixShare+"/", httpServer.ApplyCable)
	sr.HandleFunc("/"+PrefixShare+"/{id:[0-9]+}", httpServer.JoinCable)

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
