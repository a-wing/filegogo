package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"filegogo/lightcable"
	"filegogo/version"

	"github.com/gorilla/mux"
)

//go:embed dist
var dist embed.FS

func main() {

	address := flag.String("p", "0.0.0.0:8033", "set server port")
	configPath := flag.String("c", "./config.json", "use config.json")
	help := flag.Bool("h", false, "this help")
	flagVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *flagVersion {
		fmt.Printf("filegogo %s %s\n", version.Version, version.BuildTime)
		return
	}

	hub := lightcable.NewHub()
	sr := mux.NewRouter()

	sr.HandleFunc("/topic/", func(w http.ResponseWriter, r *http.Request) {
		lightcable.JoinTopic(hub, w, r)
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

	fsys, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatal(err)
	}
	sr.PathPrefix("/").Handler(http.StripPrefix("", http.FileServer(http.FS(fsys)))).Methods(http.MethodGet)

	log.Println("===============")
	log.Println("Listen Port", *address)
	log.Fatal(http.ListenAndServe(*address, sr))
}
