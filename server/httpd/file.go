package httpd

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"unicode"
)

type Meta struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	UXID string `json:"uxid"`
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func FileAttachment(w http.ResponseWriter, r *http.Request, filepath, filename string) {
	if isASCII(filename) {
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		w.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(filename))
	}
	http.ServeFile(w, r, filepath)
}

// https://stackoverflow.com/questions/53069040/checking-a-string-contains-only-ascii-characters
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
