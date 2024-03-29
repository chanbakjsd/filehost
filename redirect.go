package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// redirect checks if there's an existing .redir file and redirects the user there.
func redirect(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./hosted/" + r.URL.Path + ".redir")
	if err != nil {
		writeError(w, http.StatusNotFound)
		return
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, string(bytes), http.StatusMovedPermanently)
}
