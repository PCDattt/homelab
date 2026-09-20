package main

import (
	"fmt"
	"net/http"
	"time"
)

func sum(a, b int) int {
	return a + b
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	fmt.Println("Server started at localhost:8080")
	http.HandleFunc("/", greet)
	http.ListenAndServe(":8080", nil)
}
