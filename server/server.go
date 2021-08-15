package server

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/a-wing/lightcable"
	"github.com/gorilla/mux"
)

//go:embed dist
var dist embed.FS

const (
	Prefix = "/s"
)

func Run(address, configPath string) {
	sr := mux.NewRouter()

	cable := lightcable.New(lightcable.DefaultConfig)
	go cable.Run(context.Background())
	httpServer := NewServer(cable)

	sr.HandleFunc(Prefix+"/", httpServer.ApplyCable)
	sr.Handle(Prefix+"/{room:[0-9]+}", cable)

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
