package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Add("Cross-Origin-Embedder-Policy", "require-corp")
		fs.ServeHTTP(w, r)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: godotwebserver <web-export-dir>")
		os.Exit(1)
	}

	directory := os.Args[1]
	if ok, _ := exists(directory); !ok {
		fmt.Println("No such file or directory:", directory)
		os.Exit(1)
	}

	fmt.Println("Serving files from directory:", directory)
	fs := http.FileServer(http.Dir(directory))
	http.Handle("/", cors(fs))

	/*
	   openssl genrsa -out server.key 2048
	   openssl ecparam -genkey -name secp384r1 -out server.key
	   openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
	*/

	err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
