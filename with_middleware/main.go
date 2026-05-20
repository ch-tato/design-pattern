package main

import (
	"fmt"
	"net/http"
	"time"
)

// --- MIDDLEWARE 1: Logging ---
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("Started %s %s\n", r.Method, r.URL.Path)

		// Pass control to the next handler in the chain
		next.ServeHTTP(w, r)

		fmt.Printf("Completed in %v\n", time.Since(start))
	})
}

// --- MIDDLEWARE 2: Authentication ---
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "valid-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return // Stops the chain! The order is never placed.
		}

		// Pass control to the next handler
		next.ServeHTTP(w, r)
	})
}

// --- CORE HANDLER ---
func PlaceOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Pure business logic. No logging or auth clutter!
	fmt.Fprintln(w, "Order successfully placed for user!")
}

func main() {
	// Create the core handler
	coreHandler := http.HandlerFunc(PlaceOrderHandler)

	// Wrap it like an onion: Logger -> Auth -> Core
	secureOrderHandler := LoggingMiddleware(AuthMiddleware(coreHandler))

	http.Handle("/order", secureOrderHandler)
	http.ListenAndServe(":8080", nil)
}
