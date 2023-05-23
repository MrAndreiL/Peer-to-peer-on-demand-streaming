package utils

import (
	"net/http"

	"github.com/pkg/browser"
)

func StartStream(path, playlist string) {
	server := http.NewServeMux()
	server.Handle("/", AddHeaders(http.FileServer(http.Dir(path))))
	go http.ListenAndServe(PeerHost+VideoPort, server)
	browser.OpenURL("http://" + PeerHost + VideoPort + "/" + playlist)
}

// add CORS support.
func AddHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
