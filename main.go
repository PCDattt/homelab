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
	_, err := fmt.Fprintf(w, "Hello World! %s", time.Now())
	if err != nil {
		fmt.Println("Error :v", err)
	}
}

func main() {
	fmt.Println("Server started at localhost:8080")
	http.HandleFunc("/", greet)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error ListenAndServe %v", err)
	}
}
