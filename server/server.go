package server

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"io/ioutil"
	"strings"

	"filegogo/server/httpd"
	"filegogo/server/turnd"

	"github.com/a-wing/lightcable"
	"github.com/gorilla/mux"
	"github.com/pion/webrtc/v3"
)

//go:embed build
var dist embed.FS
var RawIndexHtml string

const (
	ApiPathConfig = "/config"
	ApiPathSignal = "/s/"
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

	sr := mux.NewRouter()

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

	fsys, err := fs.Sub(dist, "build")
	if err != nil {
		log.Fatal(err)
	}

	// read index.html file into memory
	index_, err2 := fsys.Open("index.html");
	if err2 != nil {
		log.Fatal(err2.Error());
	}
  data, _ := ioutil.ReadAll(index_);
	index_.Close();
	RawIndexHtml = string(data);
	// if exist __SUB_FOLDER__, replace it by config: SubFolder
	if strings.Contains(RawIndexHtml, "__SUB_FOLDER__") {
    RawIndexHtml = strings.ReplaceAll(RawIndexHtml, "__SUB_FOLDER__", cfg.Http.SubFolder)
	}

	//sr.PathPrefix("/").Handler(http.FileServer(httpd.NewSPA("index.html", http.FS(fsys)))).Methods(http.MethodGet)
	sr.PathPrefix("/").Handler(httpd.NewSPA(RawIndexHtml, http.FS(fsys))).Methods(http.MethodGet)

	log.Printf("=== Listen Port: %s ===\n", cfg.Http.Listen)
	log.Fatal(http.ListenAndServe(cfg.Http.Listen, sr))
}
