package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"filegogo/server/httpd"

	"github.com/gorilla/mux"
	"github.com/qingstor/go-mime"
)

func (h *Handler) NewBoxFile(w http.ResponseWriter, r *http.Request) {
	uxid := mux.Vars(r)["room"]

	f, fh, err := r.FormFile("file")
	defer f.Close()
	if err != nil {
		return
	}

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

	action := "http"
	if t := r.URL.Query().Get("action"); t != "" {
		action = t
	}

	m := &httpd.Meta{
		Name: fh.Filename,
		Size: fh.Size,
		Type: mime.DetectFileExt(strings.TrimPrefix(path.Ext(fh.Filename), ".")),
		UXID: uxid,

		Action: action,
		Remain: remain,
		Expire: expire,
	}

	h.store.Put(uxid, m)

	httpd.SaveUploadedFile(fh, path.Join(h.cfg.Http.StoragePath, uxid))

	w.Header().Add("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Println(err)
	}
}

func (h *Handler) GetBoxFile(w http.ResponseWriter, r *http.Request) {
	room := mux.Vars(r)["room"]
	var m httpd.Meta

	h.store.Get(room, &m)
	if m.UXID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: Need transaction
	m.Remain = m.Remain - 1
	h.store.Put(room, &m)

	httpd.FileAttachment(w, r, path.Join(h.cfg.Http.StoragePath, m.UXID), m.Name)
	if m.Remain == 0 {
		h.store.Del(room)
		os.Remove(path.Join(h.cfg.Http.StoragePath, m.UXID))
	}
}

func (h *Handler) DelBoxFile(w http.ResponseWriter, r *http.Request) {
	room := mux.Vars(r)["room"]
	var m httpd.Meta
	h.store.Get(room, &m)
	h.store.Del(room)
	os.Remove(path.Join(h.cfg.Http.StoragePath, m.UXID))
}

func (h *Handler) GetBoxInfo(w http.ResponseWriter, r *http.Request) {
	room := mux.Vars(r)["room"]
	var m httpd.Meta
	err := h.store.Get(room, &m)
	if time.Now().After(m.Expire) || m.Remain == 0 {
		h.store.Del(room)
		os.Remove(path.Join(h.cfg.Http.StoragePath, m.UXID))
		err = errors.New("key already expired or remain 0")
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		data, _ := json.Marshal(m)
		w.Header().Add("Content-type", "application/json")
		w.Write(data)
	}
}