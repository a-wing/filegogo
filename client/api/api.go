package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"filegogo/server/config"
	"filegogo/server/httpd"
)

const (
	ApiPathConfig = "/api/config"
	ApiPathSignal = "/api/signal/"

	ApiPathBox = "/api/box/"
	ApiPathRaw = "/api/raw/"
)

func shareGetRoom(addr string) string {
	if u, err := url.Parse(addr); err == nil {
		if arr := strings.Split(u.Path, "/"); len(arr) > 0 {
			if ok, _ := regexp.MatchString(`\d`, arr[len(arr)-1]); ok {
				return arr[len(arr)-1]
			}
		}
	}
	return ""
}

func shareGetServer(addr string) string {
	return strings.TrimSuffix(addr, "/"+shareGetRoom(addr))
}

type Api struct {
	server string
	room   string
}

func NewApi(server string) *Api {
	return &Api{
		server: shareGetServer(server),
		room:   shareGetRoom(server),
	}
}

func (a *Api) addressConfig() string {
	return a.server + ApiPathConfig
}

func (a *Api) addressSignal() string {
	return a.server + ApiPathSignal
}

func (a *Api) RoomAddress() string {
	return a.server + ApiPathSignal + a.room
}

func (a *Api) ToShare() string {
	return a.server + "/" + a.room
}

func (a *Api) GetConfig() (*config.ApiConfig, error) {
	res, err := http.Get(a.addressConfig())
	if err != nil {
		return nil, err
	}

	var cfg config.ApiConfig
	err = json.NewDecoder(res.Body).Decode(&cfg)
	return &cfg, err
}

func (a *Api) NewRoom() (string, error) {
	room, err := apiNewRoom(a.addressSignal())
	a.room = room
	return room, err
}

func (a *Api) HasRoom() bool {
	return a.room != ""
}

func apiNewRoom(addr string) (string, error) {
	res, err := http.Get(addr)
	if err != nil {
		return "", err
	}
	var msg httpd.MessageHello
	if err := json.NewDecoder(res.Body).Decode(&msg); err != nil {
		return "", err
	}
	return msg.Room, nil
}
