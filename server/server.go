package server

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"filegogo/server/httpd"
	"filegogo/server/turnd"

	"github.com/a-wing/lightcable"
	"github.com/gorilla/mux"
	"github.com/pion/webrtc/v3"
	"github.com/qingstor/go-mime"
	"github.com/rs/xid"

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

func Run(cfg *Config) {
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

	db, err := bolt.Open(path.Join(cfg.Http.StoragePath, dbName), 0600, nil)
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

	sr.HandleFunc(ApiPathConfig, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")

		var builtInICEServer *webrtc.ICEServer
		if cfg.Turn != nil {
			uaername, password := turnd.RandomUser()
			turndServer.NewUser(uaername + ":" + password)

			builtInICEServer = &webrtc.ICEServer{
				URLs:       []string{"turn:" + cfg.Turn.Listen},
				Username:   uaername,
				Credential: password,
			}
		}

		configuration := &ApiConfig{
			ICEServers: cfg.ICEServers,
		}

		if builtInICEServer != nil {
			configuration.ICEServers = append([]webrtc.ICEServer{*builtInICEServer}, cfg.ICEServers...)
		}

		if err := json.NewEncoder(w).Encode(configuration); err != nil {
			log.Println(err)
		}
	})

	sr.HandleFunc(ApiPathBoxInfo+"{room:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		room := mux.Vars(r)["room"]
		var m httpd.Meta
		err := store.Get(room, &m)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			data, _ := json.Marshal(m)
			w.Header().Add("Content-type", "application/json")
			w.Write(data)
		}
	})

	sr.HandleFunc(ApiPathBoxFile+"{room:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		uxid := xid.New().String()

		f, fh, err := r.FormFile("f")
		if err != nil {
			return
		}
		f.Close()

		remain := httpd.DefaultBoxRemain
		if t := r.URL.Query().Get("remain"); t != "" {
			remain, _ = strconv.Atoi(t)
		}

		expire := time.Now().Add(httpd.DefaultBoxExpire)
		if t := r.URL.Query().Get("expire"); t != "" {
			if tt, err := time.ParseDuration(t); err == nil {
				expire = time.Now().Add(tt)
			}
		}

		store.Put(mux.Vars(r)["room"], &httpd.Meta{
			Name: fh.Filename,
			Size: fh.Size,
			Type: mime.DetectFileExt(strings.TrimPrefix(path.Ext(fh.Filename), ".")),
			UXID: uxid,

			Remain: remain,
			Expire: expire,
		})

		httpd.SaveUploadedFile(fh, path.Join(cfg.Http.StoragePath, uxid))

	}).Methods(http.MethodPost)

	sr.HandleFunc(ApiPathBoxFile+"{room:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		room := mux.Vars(r)["room"]
		var m httpd.Meta

		store.Get(room, &m)
		if m.UXID == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// TODO: Need transaction
		m.Remain = m.Remain - 1
		store.Put(room, &m)

		httpd.FileAttachment(w, r, path.Join(cfg.Http.StoragePath, m.UXID), m.Name)
		if m.Remain == 0 {
			store.Delete(room)
			os.Remove(path.Join(cfg.Http.StoragePath, m.UXID))
		}
	}).Methods(http.MethodGet)

	sr.HandleFunc(ApiPathBoxFile+"{room:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		room := mux.Vars(r)["room"]
		var m httpd.Meta
		store.Get(room, &m)
		store.Delete(room)
		os.Remove(path.Join(cfg.Http.StoragePath, m.UXID))

	}).Methods(http.MethodDelete)

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
