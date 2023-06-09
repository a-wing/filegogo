package api

import (
	"encoding/json"
	"log"
	"net/http"

	"filegogo/server/config"
	"filegogo/server/turnd"

	"github.com/djherbis/stow/v4"
	"github.com/pion/webrtc/v3"
)

type Handler struct {
	cfg         *config.Config
	store       *stow.Store
	turndServer *turnd.Server
}

func NewHandler(cfg *config.Config, store *stow.Store, turndServer *turnd.Server) *Handler {
	return &Handler{
		cfg:         cfg,
		store:       store,
		turndServer: turndServer,
	}
}

func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	var builtInICEServer *webrtc.ICEServer
	if h.cfg.Turn != nil {
		uaername, password := turnd.RandomUser()
		h.turndServer.NewUser(uaername + ":" + password)

		builtInICEServer = &webrtc.ICEServer{
			URLs:       []string{"turn:" + h.cfg.Turn.Listen},
			Username:   uaername,
			Credential: password,
		}
	}

	configuration := &config.ApiConfig{
		ICEServers: h.cfg.ICEServers,
	}

	if builtInICEServer != nil {
		configuration.ICEServers = append([]webrtc.ICEServer{*builtInICEServer}, h.cfg.ICEServers...)
	}

	if err := json.NewEncoder(w).Encode(configuration); err != nil {
		log.Println(err)
	}
}
