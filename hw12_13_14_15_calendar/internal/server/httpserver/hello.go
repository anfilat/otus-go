package httpserver

import "net/http"

func handleHello(w http.ResponseWriter, _ *http.Request) {
	//nolint:errcheck
	w.Write([]byte("Hello, world\n"))
}
