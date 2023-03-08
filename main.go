package main

import (
	"flag"
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
	var directory string
	flag.StringVar(&directory, "directory", ".", "path to godot web export directory")

	if ok, _ := exists(directory); !ok {
		os.Exit(1)
	}

	fmt.Println("Using directory:", directory)

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
