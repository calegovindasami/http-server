package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)


func main() {
	http.HandleFunc("/", root)

	err := http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server has been shutdown\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}


func root(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request: GET /\n")
	io.WriteString(w, "This is my server root\n")
}