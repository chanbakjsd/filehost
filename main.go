package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	err := os.Mkdir("./hosted", 0755)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/upload", upload)
	http.Handle("/hosted/", http.StripPrefix("/hosted/", http.FileServer(http.Dir("./hosted"))))
	http.Handle("/", http.FileServer(http.Dir("./static")))

	go cleanFolder()

	fmt.Println("Server initialized.")

	if err := http.ListenAndServe(host, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
