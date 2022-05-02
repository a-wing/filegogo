package api

import (
	"testing"
)

func TestShareGetRoom(t *testing.T) {
	url1, r1 := "http://localhost:8080", ""
	url2, r2 := "http://localhost:8080/api", ""
	url3, r3 := "http://localhost:8080/api/1234", "1234"
	url4, r4 := "http://localhost:8080/api/qwq/1234", "1234"
	if shareGetRoom(url1) != r1 {
		t.Error("Should equal")
	}

	if shareGetRoom(url2) != r2 {
		t.Error("Should equal")
	}

	if shareGetRoom(url3) != r3 {
		t.Error("Should equal")
	}

	if shareGetRoom(url4) != r4 {
		t.Error("Should equal")
	}
}

func TestShareGetServer(t *testing.T) {
	url1, r1 := "http://localhost:8080", "http://localhost:8080"
	url2, r2 := "http://localhost:8080/api", "http://localhost:8080/api"
	url3, r3 := "http://localhost:8080/1204", "http://localhost:8080"
	url4, r4 := "http://localhost:8080/qwq/1234", "http://localhost:8080/qwq"
	if shareGetServer(url1) != r1 {
		t.Error("Should equal")
	}

	if shareGetServer(url2) != r2 {
		t.Error("Should equal")
	}

	if shareGetServer(url3) != r3 {
		t.Error("Should equal")
	}

	if shareGetServer(url4) != r4 {
		t.Error("Should equal")
	}
}
