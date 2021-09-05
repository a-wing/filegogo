package share

import (
	"filegogo/server"
	"testing"
)

func TestIsShareInit(t *testing.T) {
	initAddr := "http://localhost:8080"
	if IsShareInit(initAddr) {
		t.Error("Not init")
	}

	joinAddr := "http://localhost:8080/1234"
	if !IsShareInit(joinAddr) {
		t.Error("Need init")
	}
}

func TestShareToLinks(t *testing.T) {
	req := "http://localhost:8033/1024"
	res := "http://localhost:8033" + server.Prefix + "1024"

	if ShareToLinks(req) != res {
		t.Error("Should equal")
	}
}

func TestLinksToShare(t *testing.T) {
	req := "http://localhost:8033" + server.Prefix + "1024"
	res := "http://localhost:8033/1024"

	if LinksToShare(req) != res {
		t.Error("Should equal")
	}
}
