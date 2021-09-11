package server

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"net/url"

	"filegogo/server/turnd"

	"github.com/a-wing/lightcable"
	"github.com/gorilla/mux"
)

//go:embed build
var dist embed.FS

const (
	Prefix = "/s/"
)

// Fork From the: https://pkg.go.dev/net/http#StripPrefix
func NoPrefix(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			h.ServeHTTP(w, r)
		} else {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = ""
			r2.URL.RawPath = ""
			h.ServeHTTP(w, r2)
		}
	})
}

func Run(cfg *Config) {
	turndServer := turnd.New(&turnd.Config{
		Username:     "filegogo",
		Password:     "filegogo",
		Realm:        "filegogo",
		Listen:       "0.0.0.0:3478",
		PublicIP:     "0.0.0.0",
		RelayMinPort: 49160,
		RelayMaxPort: 49200,
	})
	turndServer.NewUser("filegogo")
	turnSrv, err := turndServer.Run()
	if err != nil {
		panic(err)
	}
	defer turnSrv.Close()

	sr := mux.NewRouter()

	cable := lightcable.New(lightcable.DefaultConfig)
	go cable.Run(context.Background())
	httpServer := NewServer(cable)

	sr.HandleFunc(Prefix, httpServer.ApplyCable)
	sr.Handle(Prefix+"{room:[0-9]+}", cable)

	sr.HandleFunc("/config.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		if err := json.NewEncoder(w).Encode(cfg.IcsServers); err != nil {
			log.Println(err)
		}
	})

	fsys, err := fs.Sub(dist, "build")
	if err != nil {
		log.Fatal(err)
	}

	sr.PathPrefix("/{id:[0-9]+}").Handler(NoPrefix(http.FileServer(http.FS(fsys)))).Methods(http.MethodGet)
	sr.PathPrefix("/").Handler(http.FileServer(http.FS(fsys))).Methods(http.MethodGet)

	log.Printf("=== Listen Port: %s ===\n", cfg.Server)
	log.Fatal(http.ListenAndServe(cfg.Server, sr))
}
