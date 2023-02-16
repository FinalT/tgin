package tgin

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	var s Server
	http.ListenAndServe(":8080", s)

	http.ListenAndServeTLS(":4000",
		"cret file", "key file", s)

	s.Start(":8081")
}
