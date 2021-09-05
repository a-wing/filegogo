package server

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"

	"github.com/a-wing/lightcable"
	"github.com/gorilla/mux"
)

//go:embed dist
var dist embed.FS

const (
	Prefix = "/s"
)

func Run(cfg *Config) {
	sr := mux.NewRouter()

	cable := lightcable.New(lightcable.DefaultConfig)
	go cable.Run(context.Background())
	httpServer := NewServer(cable)

	sr.HandleFunc(Prefix+"/", httpServer.ApplyCable)
	sr.Handle(Prefix+"/{room:[0-9]+}", cable)

	sr.HandleFunc("/config.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		if err := json.NewEncoder(w).Encode(cfg.IcsServers); err != nil {
			log.Println(err)
		}
	})

	fsys, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatal(err)
	}
	sr.PathPrefix("/").Handler(http.StripPrefix("", http.FileServer(http.FS(fsys)))).Methods(http.MethodGet)

	log.Printf("=== Listen Port: %s ===\n", cfg.Server)
	log.Fatal(http.ListenAndServe(cfg.Server, sr))
}
