package pool

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func getFileContentType(file *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	c, err := file.Read(buffer)
	defer file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer[:c])
	return contentType, nil
}

func HelpCreateTmpFile(name string, size int) (string, string, error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	sum := md5.New()

	chunkSize := 4096

	count := size / chunkSize
	remain := size % chunkSize

	// Default /tmp
	f, err := os.CreateTemp("", name)
	if err != nil {
		return "", "", err
	}
	// count
	data := make([]byte, 4096)
	for i := 0; i < count; i++ {
		if _, err := r.Read(data); err != nil {
			return "", "", err
		}

		sum.Write(data)

		if _, err := f.Write(data); err != nil {
			return "", "", err
		}
	}

	// remain
	if _, err := r.Read(data); err != nil {
		return "", "", err
	}

	sum.Write(data[:remain])

	if _, err := f.Write(data[:remain]); err != nil {
		return "", "", err
	}

	if err := f.Close(); err != nil {
		return "", "", err
	}
	return f.Name(), fmt.Sprintf("%x", sum.Sum(nil)), nil
}
