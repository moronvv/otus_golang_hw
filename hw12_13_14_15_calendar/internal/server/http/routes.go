package internalhttp

import (
	"io"
	"net/http"
)

func helloRoute(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello World!\n")
}
