package server

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"filegogo/server/api"
	"filegogo/server/config"
	"filegogo/server/httpd"
	"filegogo/server/turnd"

	"github.com/a-wing/lightcable"
	"github.com/gorilla/mux"

	"github.com/djherbis/stow/v4"
	bolt "go.etcd.io/bbolt"
)

//go:embed build
var dist embed.FS

const (
	ApiPathConfig = "/api/config"
	ApiPathSignal = "/api/signal/"

	ApiPathBoxInfo = "/api/info/"
	ApiPathBoxFile = "/api/file/"

	dbName = "store.db"
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
	if err := os.MkdirAll(cfg.Http.StoragePath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(path.Join(cfg.Http.StoragePath, dbName), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	store := stow.NewJSONStore(db, []byte("room"))

	sr := mux.NewRouter().PathPrefix("/"+cfg.Http.PathPrefix).Subrouter()

	cable := lightcable.New(lightcable.DefaultConfig)
	go cable.Run(context.Background())
	httpServer := httpd.NewServer(cable, cfg.Http)

	sr.HandleFunc(ApiPathSignal, httpServer.ApplyCable)
	sr.Handle(ApiPathSignal+"{room:[0-9]+}", cable)

	hander := api.NewHandler(cfg, store, turndServer)

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

	go run(context.Background(), func() {
		now := time.Now()
		store.ForEach(func(key string, val httpd.Meta) {
			if now.After(val.Expire) {
				store.Delete(key)
				os.Remove(path.Join(cfg.Http.StoragePath, val.UXID))
			}
		})
	})

	log.Printf("=== Listen Port: %s ===\n", cfg.Http.Listen)
	log.Fatal(http.ListenAndServe(cfg.Http.Listen, sr))
}
