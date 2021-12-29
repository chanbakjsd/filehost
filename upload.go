package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// upload handles the upload endpoint.
func upload(w http.ResponseWriter, r *http.Request) {
	// Checks if it's a valid request or has reached request limit.
	if r.Method != "POST" {
		writeError(w, http.StatusBadRequest)
		return
	}
	if hasHitRequestLimit(r.RemoteAddr) {
		writeError(w, http.StatusTooManyRequests)
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

	// Sanitize file name and attempt to save the file.
	resultFileName := sanitizeFileName(files[0].Filename)
	err := saveFile(resultFileName, files[0])
	if err != nil {
		fmt.Println(err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	// Return that it has been uploaded successfully.
	w.Write([]byte(domain + "/hosted/" + resultFileName))
}

// sanitizeFileName assigns a new filename to the file with a random name and
// attempts to keep the extension provided that it's lowercase ASCII and not
// redir (which is reserved for redirection).
func sanitizeFileName(filename string) string {
	filenameSplit := strings.Split(filename, ".")
	ext := ""
	if len(filenameSplit) > 1 {
		final := filenameSplit[len(filenameSplit)-1]
		safe := true
		for _, v := range []byte(final) {
			if (v < 'a' || v > 'z') && (v < '0' || v > '9') {
				safe = false
				break
			}
		}
		if safe && final != "redir" { // Redir is reserved for redirection.
			ext = "." + final
		}
	}
	return strconv.FormatInt(getNextFileNum(), 36) + ext
}

// saveFile saves the provided multipart file as the provided name in the hosted
// directory.
func saveFile(targetName string, sourceFile *multipart.FileHeader) error {
	// Create the file to write to.
	file, err := os.Create("./hosted/" + targetName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Prepare the source.
	source, err := sourceFile.Open()
	if err != nil {
		return err
	}

	// And save the file.
	_, err = io.Copy(file, source)
	return err
}
