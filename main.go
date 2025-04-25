package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Book struct {
	ID int `json:"id,omitempty"`
	Title string `json:"title"`
	Author string `json:"author"`
}

var libraryCache = make(map[int]Book)

var cacheMutex sync.RWMutex

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", root)
	mux.HandleFunc("POST /book", createBook)
	mux.HandleFunc("DELETE /book/{id}", deleteBook)
	mux.HandleFunc("GET  /book", getBooks)
	mux.HandleFunc("GET /book/{id}", getBookById)

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

func createBook(w http.ResponseWriter, r*http.Request) {
	fmt.Printf("POST /book\n")

	var book  Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if book.Title == "" || book.Author == "" {
		http.Error(w, "title and author is required", http.StatusBadRequest)
		return
	}


	id := len(libraryCache) + 1
	cacheMutex.Lock()
	libraryCache[len(libraryCache) + 1] = book
	cacheMutex.Unlock()
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, id)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DELETE /book\n")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheMutex.RLock()
	_, ok := libraryCache[id]
	cacheMutex.RUnlock()

	if !ok {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	delete(libraryCache, id)
	w.WriteHeader(http.StatusNoContent)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET /book\n")

	books := make([]Book, 0, len(libraryCache))


	for id, book := range libraryCache {
		bookCopy := book
		bookCopy.ID = id
		books = append(books, bookCopy)
	}

	json, err := json.Marshal(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func getBookById(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET /book/{id}\n")


	 id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheMutex.RLock()
	book, ok := libraryCache[id]
	cacheMutex.RUnlock()

	if !ok {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	json, err := json.Marshal(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

