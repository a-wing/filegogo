package api

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"filegogo/server/httpd"
	"filegogo/server/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/qingstor/go-mime"
)

func (h *Handler) NewBox(w http.ResponseWriter, r *http.Request) {
	uxid := h.genID()
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

	m := &httpd.Box{
		UXID: uxid,

		Secret: utils.GenSecret(16),
		Action: action,
		Remain: remain,
		Expire: expire,
	}

	if render.GetRequestContentType(r) == render.ContentTypeJSON {
		render.DecodeJSON(r.Body, m)
	}

	f, fh, err := r.FormFile("file")
	if err == nil {
		defer f.Close()
		httpd.SaveUploadedFile(fh, path.Join(h.cfg.Http.StoragePath, uxid))
		m.Name = fh.Filename
		m.Size = fh.Size
		m.Type = mime.DetectFileExt(strings.TrimPrefix(path.Ext(fh.Filename), "."))
	}

	h.store.Put(uxid, m)

	render.JSON(w, r, m)
}


func (h *Handler) GetBox(w http.ResponseWriter, r *http.Request) {
	uxid := chi.URLParam(r, "ID")
	var m httpd.Box
	err := h.store.Get(uxid, &m)
	if time.Now().After(m.Expire) || m.Remain == 0 {
		h.store.Del(uxid)
		os.Remove(path.Join(h.cfg.Http.StoragePath, m.UXID))
		err = errors.New("key already expired or remain 0")
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		m.Secret = ""
		render.JSON(w, r, m)
	}
}

func (h *Handler) DelBox(w http.ResponseWriter, r *http.Request) {
	uxid := chi.URLParam(r, "ID")
	var m httpd.Box
	h.store.Get(uxid, &m)

	if m.Secret == r.URL.Query().Get("secret") {
		h.store.Del(uxid)
		os.Remove(path.Join(h.cfg.Http.StoragePath, m.UXID))
		render.NoContent(w, r)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func (h *Handler) GetRaw(w http.ResponseWriter, r *http.Request) {
	uxid := chi.URLParam(r, "ID")
	var m httpd.Box

	h.store.Get(uxid, &m)
	if m.UXID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: Need transaction
	m.Remain = m.Remain - 1
	h.store.Put(uxid, &m)

	httpd.FileAttachment(w, r, path.Join(h.cfg.Http.StoragePath, m.UXID), m.Name)
	if m.Remain == 0 {
		h.store.Del(uxid)
		os.Remove(path.Join(h.cfg.Http.StoragePath, m.UXID))
	}
}
