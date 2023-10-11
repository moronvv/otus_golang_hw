package internalhttproutes

import (
	"io"
	"net/http"
)

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "pong\n")
}
