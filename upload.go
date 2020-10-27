package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//upload handles the upload endpoint.
func upload(w http.ResponseWriter, r *http.Request) {
	// Checks if it's a valid request or has reached request limit.
	if r.Method != "POST" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	if hasHitRequestLimit(r.RemoteAddr) {
		http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
		return
	}
	if err := r.ParseMultipartForm(maxMemoryPerRequest); err != nil {
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	// Check if password is correct if mandated.
	if password != "" && (len(r.MultipartForm.Value["pass"]) != 1 || r.MultipartForm.Value["pass"][0] != password) {
		http.Error(w, "Password protected instance", http.StatusBadRequest)
		return
	}

	// Check that there are exactly one file being uploaded.
	files := r.MultipartForm.File["file"]
	if len(files) != 1 {
		http.Error(w, "Expected exactly one file", http.StatusBadRequest)
		return
	}

	// Check that the file is not too big.
	if hasHitSizeLimit(r.RemoteAddr, files[0].Size) {
		http.Error(w, "File too big", http.StatusTooManyRequests)
		return
	}

	// "Sanitize" the name. Give it a random name and
	// reuse the extension if it's lowercase letters to make it easier for users.
	filenameSplit := strings.Split(files[0].Filename, ".")
	ext := ""
	if len(filenameSplit) >= 2 {
		final := filenameSplit[len(filenameSplit)-1]
		safe := true
		for _, v := range []byte(final) {
			if v < 'a' || v > 'z' {
				safe = false
			}
		}
		if safe {
			ext = "." + filenameSplit[len(filenameSplit)-1]
		}
	}

	// Create the file to write to.
	resultFileName := strconv.FormatInt(getNextFileNum(), 36) + ext
	file, err := os.Create("./hosted/" + resultFileName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Prepare the source.
	source, err := files[0].Open()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// And save the file.
	_, err = io.Copy(file, source)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Return that it has been uploaded successfully.
	w.Write([]byte(domain + "/hosted/" + resultFileName))
}
