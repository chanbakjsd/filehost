package main

import (
	"net/http"
)

// writeError writes the provided error to the ResponseWriter.
func writeError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
