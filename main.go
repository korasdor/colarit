package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
)

func main() {

	port := os.Getenv("PORT")


	if port == "" {
		port = "8080"
	}

	fmt.Printf("listening on %s...", ":"+port)
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	mux.Vars(req)

	fmt.Fprintf(res, "hello, world from %s", runtime.Version())
}
