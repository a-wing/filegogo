package httpd

import (
	"net/http"
)

func NewSPA(path string, fs http.FileSystem) *SPA {
	return &SPA{
		path: path,
		fs:   fs,
	}
}

type SPA struct {
	path string
	fs   http.FileSystem
}

func (s *SPA) Open(name string) (http.File, error) {
	if file, err := s.fs.Open(name); err != nil {
		return s.fs.Open(s.path)
	} else {
		return file, err
	}
}
