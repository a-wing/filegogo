package api

import (
	"encoding/json"
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
	//"github.com/rs/xid"
)

func (h *Handler) NewBoxFile(w http.ResponseWriter, r *http.Request) {
	//uxid := xid.New().String()
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

	m := &httpd.Meta{
		Name: fh.Filename,
		Size: fh.Size,
		Type: mime.DetectFileExt(strings.TrimPrefix(path.Ext(fh.Filename), ".")),
		UXID: uxid,

		Remain: remain,
		Expire: expire,
	}

	//h.store.Put(mux.Vars(r)["room"], m)
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
		h.store.Delete(room)
		os.Remove(path.Join(h.cfg.Http.StoragePath, m.UXID))
	}
}

func (h *Handler) DelBoxFile(w http.ResponseWriter, r *http.Request) {
	room := mux.Vars(r)["room"]
	var m httpd.Meta
	h.store.Get(room, &m)
	h.store.Delete(room)
	os.Remove(path.Join(h.cfg.Http.StoragePath, m.UXID))
}

func (h *Handler) GetBoxInfo(w http.ResponseWriter, r *http.Request) {
	room := mux.Vars(r)["room"]
	var m httpd.Meta
	err := h.store.Get(room, &m)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		data, _ := json.Marshal(m)
		w.Header().Add("Content-type", "application/json")
		w.Write(data)
	}
}
