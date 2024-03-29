package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// redirectRegister handles requests to create a new redirect.
func redirectRegister(w http.ResponseWriter, r *http.Request) {
	// Checks if it's a valid request or has reached request limit.
	if r.Method != "POST" {
		writeError(w, http.StatusBadRequest)
		return
	}
	if hasHitRequestLimit(r.RemoteAddr) {
		writeError(w, http.StatusTooManyRequests)
		return
	}

	err := r.ParseForm()
	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	// Check if password is correct if mandated.
	if password != "" && (len(r.Form["pass"]) != 1 || r.Form["pass"][0] != password) {
		http.Error(w, "Password protected instance", http.StatusBadRequest)
		return
	}
	// Check that there are exactly one URL being sent.
	if len(r.Form["url"]) != 1 {
		http.Error(w, "Expected exactly one target URL", http.StatusBadRequest)
		return
	}
	targetURL := r.Form["url"][0]

	// Check that the URL is not too long.
	if len(targetURL) > 8000 || hasHitSizeLimit(r.RemoteAddr, int64(len(targetURL))) {
		http.Error(w, "URL too long.", http.StatusTooManyRequests)
		return
	}

	// Make sure it's valid URL.
	parsedURL, err := url.Parse(targetURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		http.Error(w, "Invalid URL provided.", http.StatusBadRequest)
		return
	}

	// Create the file to write to.
	resultID := strconv.FormatInt(getNextFileNum(), 36)
	file, err := os.Create("./hosted/" + resultID + ".redir")
	if err != nil {
		fmt.Println(err)
		writeError(w, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Save the URL into the file.
	_, err = file.Write([]byte(targetURL))
	if err != nil {
		fmt.Println(err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(domain + "/r/" + resultID))
}
