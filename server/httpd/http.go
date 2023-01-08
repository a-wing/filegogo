package httpd

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

func NewSPA(raw_index string, fs http.FileSystem) *SPA {
	return &SPA{
		RawIndex: raw_index,
		fs:   fs,
	}
}

type SPA struct {
	RawIndex string
	fs   http.FileSystem
}

func (s* SPA) _index_serve(w http.ResponseWriter) {
	w.Header().Add("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s.RawIndex))
}

func (s* SPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// root or index.html
	if path == "/" || path == "/index.html" {
		// serve index.html
		s._index_serve(w)
		log.Info("response index: index.html")
		return
	} else {
  	// check whether a file exists at the given path
	  file, err := s.fs.Open(path)
		if err != nil {
			// file does not exist, serve index.html
			s._index_serve(w)
			log.Info("response not found: index.html")
			return
		}
		defer file.Close()
	}

	// otherwise, use http.FileServer to serve the embed dir
	log.Info("response: ", path)
	http.FileServer(s.fs).ServeHTTP(w, r)
}
