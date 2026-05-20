package main

import (
	"fmt"
	"net/http"
	"time"
)

// One massive handler doing everything
func PlaceOrderHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Logging Logic (Duplicated everywhere)
	start := time.Now()
	fmt.Printf("Started %s %s\n", r.Method, r.URL.Path)

	// 2. Authentication Logic (Duplicated everywhere)
	token := r.Header.Get("Authorization")
	if token != "valid-token" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 3. CORE BUSINESS LOGIC (The only thing this function should care about)
	fmt.Fprintln(w, "Order successfully placed for user!")

	// Logging continued...
	fmt.Printf("Completed in %v\n", time.Since(start))
}

func main() {
	http.HandleFunc("/order", PlaceOrderHandler)
	http.ListenAndServe(":8080", nil)
}
