package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Person struct {
	Name string `json:"name"`
	LastName string `json:"lastName"`
}


func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", root)
	mux.HandleFunc("POST /hello", hello)

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
	fmt.Fprintf(w, "This is my server root\n")
}

func hello(w http.ResponseWriter, r *http.Request) {
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("POST /hello\n")
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello %s %s!", person.Name, person.LastName)
}