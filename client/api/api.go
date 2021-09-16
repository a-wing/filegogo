package api

import (
	"encoding/json"
	"net/http"

	"filegogo/server/httpd"
)

func GetRoom(addr string) (string, error) {
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
