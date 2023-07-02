package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"

	"filegogo/server/config"
	"filegogo/server/httpd"
	"filegogo/server/store"
)

func TestBox(t *testing.T) {
	d1 := []byte("hello world!")
	f, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Write(d1)

	f2, err := os.Open(f.Name())
	if err != nil {
		t.Error(err)
	}
	stat, err := f2.Stat()
	if err != nil {
		t.Error(err)
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("file", stat.Name())
	if err != nil {
		t.Error(err)
	}

	b, err := ioutil.ReadAll(f2)
	if err != nil {
		t.Error(err)
	}

	part.Write(b)
	writer.Close()

	cfg := &config.Config{
		Http: &httpd.Config{
			StoragePath: os.TempDir(),
		},
	}

	handler := NewHandler(cfg, store.NewStore(), nil)

	req := httptest.NewRequest("POST", "http://example.com/foo", buf)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	handler.NewBox(w, req)

	m := &httpd.Box{}
	if err := json.NewDecoder(w.Body).Decode(m); err != nil {
		t.Error(err)
	}
	if m.Size != int64(len(d1)) {
		t.Error(m.Size)
	}

}
