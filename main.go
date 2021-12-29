package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func main() {
	err := os.Mkdir("./hosted", 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(err)
	}

	files := http.StripPrefix("/hosted/", http.FileServer(http.Dir("./hosted")))

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/shorten", redirectRegister)
	http.HandleFunc("/hosted/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/hosted/" {
			http.NotFound(w, r)
			return
		}
		files.ServeHTTP(w, r)
	})
	http.Handle("/r/", http.StripPrefix("/r/", http.HandlerFunc(redirect)))
	http.Handle("/", http.FileServer(http.Dir("./static")))

	go cleanFolder()

	fmt.Println("Server initialized.")

	if err := http.ListenAndServe(host, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
