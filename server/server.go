package server

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"filegogo/server/api"
	"filegogo/server/config"
	"filegogo/server/httpd"
	"filegogo/server/store"
	"filegogo/server/turnd"

	"github.com/a-wing/lightcable"
	"github.com/gorilla/mux"
)

//go:embed build
var dist embed.FS

const (
	ApiPathConfig = "/api/config"
	ApiPathSignal = "/api/signal/"

	ApiPathBoxInfo = "/api/info/"
	ApiPathBoxFile = "/api/file/"
)

func Run(cfg *config.Config) {
	var turndServer *turnd.Server
	if cfg.Turn != nil {
		log.Println("Enabled Built-in Stun And Turn Server")
		turndServer = turnd.New(cfg.Turn)
		turnSrv, err := turndServer.Run()
		if err != nil {
			panic(err)
		}
		defer turnSrv.Close()
	}

	if err := os.RemoveAll(cfg.Http.StoragePath); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(cfg.Http.StoragePath, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	sr := mux.NewRouter().PathPrefix("/"+cfg.Http.PathPrefix).Subrouter()

	cable := lightcable.New(lightcable.DefaultConfig)
	go cable.Run(context.Background())
	httpServer := httpd.NewServer(cable, cfg.Http)

	sr.HandleFunc(ApiPathSignal, httpServer.ApplyCable)
	sr.Handle(ApiPathSignal+"{room:[0-9]+}", cable)

	hander := api.NewHandler(cfg, store.NewStore(), turndServer)

	sr.HandleFunc(ApiPathConfig, hander.GetConfig)
	sr.HandleFunc(ApiPathBoxInfo+"{room:[0-9]+}", hander.GetBoxInfo)
	sr.HandleFunc(ApiPathBoxFile+"{room:[0-9]+}", hander.NewBoxFile).Methods(http.MethodPost)
	sr.HandleFunc(ApiPathBoxFile+"{room:[0-9]+}", hander.GetBoxFile).Methods(http.MethodGet)
	sr.HandleFunc(ApiPathBoxFile+"{room:[0-9]+}", hander.DelBoxFile).Methods(http.MethodDelete)

	fsys, err := fs.Sub(dist, "build")
	if err != nil {
		log.Fatal(err)
	}

	sr.PathPrefix("/").Handler(http.StripPrefix("/"+cfg.Http.PathPrefix, http.FileServer(httpd.NewSPA("index.html", http.FS(fsys))))).Methods(http.MethodGet)

	log.Printf("=== Listen Port: %s ===\n", cfg.Http.Listen)
	log.Fatal(http.ListenAndServe(cfg.Http.Listen, sr))
}
