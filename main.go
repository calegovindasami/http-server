package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)


func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", root)
	mux.HandleFunc("/hello", helloWorld)

	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server has been shutdown\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}


func root(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET /\n")
	io.WriteString(w, "This is my server root\n")
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET /hello\n")
	io.WriteString(w, "Hello World!")
}