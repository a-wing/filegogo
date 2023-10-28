package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/a-wing/lightcable"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"filegogo/server/api"
	"filegogo/server/config"
	"filegogo/server/httpd"
	"filegogo/server/store"
	"filegogo/server/turnd"
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
	cable := lightcable.New(lightcable.DefaultConfig)
	go cable.Run(context.Background())

	hander := api.NewHandler(cfg, store.NewStore(), turndServer)

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/"+cfg.Http.PathPrefix, func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Get("/config", hander.GetConfig)
			r.Handle("/signal/*", cable)

			r.Route("/raw", func(r chi.Router) {
				r.Get("/{ID}", hander.GetRaw)
			})

			r.Route("/box", func(r chi.Router) {
				r.Post("/", hander.NewBox)
				r.Get("/{ID}", hander.GetBox)
				r.Delete("/{ID}", hander.DelBox)
			})
		})

		r.Handle("/*", http.StripPrefix("/"+cfg.Http.PathPrefix, http.FileServer(httpd.NewSPA("index.html", http.FS(dist)))))
	})

	log.Printf("=== Listen Port: %s ===\n", cfg.Http.Listen)
	log.Fatal(http.ListenAndServe(cfg.Http.Listen, r))
}
