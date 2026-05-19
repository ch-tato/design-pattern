Here is the complete breakdown for the **Middleware** pattern, continuing with our **Intelligent Food Delivery Platform** study case.

### **Name**

**Middleware Pattern** (A web-specific application of the **Decorator** and **Chain of Responsibility** patterns)

---

### **Intent / Problem**

**The Problem:** In our food delivery platform, when a user sends an HTTP request to place a `DeliveryOrder`, we need to do several things before actually saving the order to the database: we need to log the incoming request, check if the user is authenticated (valid token), and ensure they aren't spamming the server (rate limiting). If we write all this logic directly inside the "Place Order" function, that function becomes massive. Worse, when we add a "Cancel Order" function, we have to copy and paste the exact same logging and authentication code.

**The Intent:**
Middleware extracts "cross-cutting concerns" (like logging, authentication, and metrics) out of your core business logic. It allows you to wrap your core HTTP handlers in layers (like an onion). The request passes through the layers to get in, and the response passes back through the layers to get out.

---

### **Solution**

**Static Structure:**

* **Core Handler:** The actual business logic function (e.g., `PlaceOrderHandler`).
* **Middleware Function:** A function that takes a Handler as an input, wraps it with some extra logic, and returns a new, wrapped Handler.

**Dynamic Behavior:**

1. A web request arrives.
2. It hits the outermost middleware (e.g., `Logger`). The logger records the time, then passes the request to the next layer.
3. It hits the next middleware (e.g., `Authenticator`). If the user is missing a token, this layer immediately rejects the request and returns a 401 Unauthorized error. The core handler is never reached.
4. If authentication passes, it reaches the `Core Handler`, which processes the food order.

---

### **Sample Code (Go)**

In Go, Middleware is perfectly natively supported by the standard `net/http` library.

#### **1. Without Middleware (The "Fat Handler" Problem)**

Here, all the logging and security logic is tangled up with the business logic.

```go
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

```

#### **2. With Middleware (The Clean, Layered Solution)**

Here, we extract logging and authentication into reusable wrappers. The core handler only deals with food orders.

```go
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

```

---

### **Discussion**

**Implementation Issues & Considerations:**

* **Order of Execution:** The order in which you chain middleware is critical. If you put `AuthMiddleware` *before* `LoggingMiddleware`, you won't have any logs for unauthorized users. You usually want Logging first, then Security/Auth, then Validation, then the Core logic.
* **Context Passing:** Often, middleware needs to pass data to the core handler (e.g., the `AuthMiddleware` extracts the User ID from the token). In Go, this is done using the `context` package attached to the `http.Request`.
* **Third-Party Routers:** While chaining by hand (`Log(Auth(Handler))`) is fine for small apps, most modern Go frameworks (like Chi, Gin, or Echo) provide a `.Use()` method to apply middleware globally across multiple routes at once, keeping the code incredibly clean.

---

### **Related Patterns**

* **Decorator:** Middleware is essentially the structural Decorator pattern applied specifically to HTTP requests. It attaches new behaviors to objects (handlers) by placing them inside special wrapper objects.
* **Chain of Responsibility:** Very similar in concept. Both pass a request along a chain of handlers. However, in standard Chain of Responsibility, a handler usually *consumes* the request and stops the chain. In Middleware, the handler usually does its job and *passes* it to the next layer, catching the response on the way back out.
* **Proxy:** A Proxy controls access to an object. Middleware (like the `AuthMiddleware`) acts as a protective proxy for the core business logic.