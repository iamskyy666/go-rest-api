# SECURITY-HEADERS Middleware() 🛡️✨

```go
package middlewares

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-DNS-Prefetch-Control","off")	
	w.Header().Set("X-Frame-Options","DENY")
	w.Header().Set("X-XSS-Protection","1;mode-block")
	w.Header().Set("X-Content-Type-Options","nosniff")
	w.Header().Set("Strict Transport Security","max-age=63072000;includeSubDomains;preload")
	w.Header().Set("Content-Security-Policy","default-src 'self'")
	w.Header().Set("Referrer-Policy","no-referrer")
	next.ServeHTTP(w,r)
	})
}

//! BASIC SNIPPET:
// 📍 root/internal/api/middlewares/security_headers.middleware.go

//import "net/http"
// func securityHeaders(next http.Handler) http.Handler{
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	next.ServeHTTP(w,r)
// 	})
// }
```
---

### 🧩 What This Code Is Doing

This function defines a **middleware** in Go (for our REST API).
A **middleware** is a piece of code that runs **before** (or sometimes after) our main handler executes.

It usually does things like:

* Add headers
* Authenticate users
* Log requests
* Handle CORS
* Sanitize input
  etc.

---

### ✅ Full Breakdown

#### 📦 1. Package declaration

```go
package middlewares
```

This means this file belongs to a package named **`middlewares`**, and other parts of our app can import it like:

```go
import "ourapp/internal/api/middlewares"
```

---

#### 📥 2. Import

```go
import "net/http"
```

We need the `net/http` package because middleware in Go works with **HTTP handlers**, which come from this package.

---

#### ⚙️ 3. Function definition

```go
func SecurityHeaders(next http.Handler) http.Handler {
```

* `SecurityHeaders` is the function name (exported because it starts with a capital letter).
* It **takes** one parameter:
  `next http.Handler` → this is the *next handler* in the chain (e.g., our main endpoint handler).
* It **returns** another `http.Handler`.
  So this function wraps around another handler — that’s why it’s called middleware.

---

#### 🚀 4. Return a new handler

```go
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
```

* We’re returning a new handler created with `http.HandlerFunc`.
* It takes a function with the signature `(w http.ResponseWriter, r *http.Request)`, which is what all Go HTTP handlers look like.

---

#### 🧱 5. Setting security headers

Inside the returned function, we set some HTTP **security headers**:

```go
w.Header().Set("X-DNS-Prefetch-Control", "off")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("X-XSS-Protection", "1;mode-block")
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("Strict-Transport-Security", "max-age=63072000;includeSubDomains;preload")
w.Header().Set("Content-Security-Policy", "default-src 'self'")
w.Header().Set("Referrer-Policy", "no-referrer")
```

Here’s what each one does 👇

| Header                                          | Purpose                                                                                   |
| ----------------------------------------------- | ----------------------------------------------------------------------------------------- |
| **X-DNS-Prefetch-Control**                      | Controls browser DNS prefetching. “off” means don’t prefetch — small privacy improvement. |
| **X-Frame-Options: DENY**                       | Prevents our site from being loaded inside an `<iframe>` (helps against clickjacking).   |
| **X-XSS-Protection: 1; mode=block**             | Enables browser’s XSS protection (older browsers).                                        |
| **X-Content-Type-Options: nosniff**             | Prevents MIME type sniffing — the browser will trust our declared `Content-Type`.        |
| **Strict-Transport-Security**                   | Forces browsers to use HTTPS for all future requests to this domain.                      |
| **Content-Security-Policy: default-src 'self'** | Only allow resources (scripts, images, etc.) from our own domain.                        |
| **Referrer-Policy: no-referrer**                | Don’t send referrer information to other sites.                                           |

These are **important for web security** — they make our API safer by preventing common attacks like:

* XSS (Cross-Site Scripting)
* Clickjacking
* Information leakage

---

#### 🔁 6. Call the next handler

```go
next.ServeHTTP(w, r)
```

After setting the headers, we pass control to the next handler in the chain (like our main route handler).
If you forget this line, the request would stop here and never reach our actual API logic!

---

#### 🧱 BASIC SNIPPET (for reference)

The basic snippet at the bottom of our code:

```go
func securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)
    })
}
```

This is a **skeleton** for a middleware — it doesn’t do anything yet, but shows the structure all middlewares in Go follow.

---

### 🧠 TL;DR (In Simple Words)

* Middleware = wrapper around our handler.
* This one adds HTTP security headers to every response.
* After adding headers, it calls the next handler so the request continues.

---

### 💡 Example of Using It

In our main `router.go` or `main.go`:

```go
mux := http.NewServeMux()

mux.Handle("/api/v1/data", middlewares.SecurityHeaders(http.HandlerFunc(myHandler)))

http.ListenAndServe(":8080", mux)
```

Now every response from `/api/v1/data` will include those security headers.
---

# RATE-LIMITER Middleware() ⚠️📈

```go
package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu  sync.Mutex
	visitors  map[string]int
	limit  int
	resetTime  time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration)*rateLimiter{
	rl:= &rateLimiter{
		visitors: make(map[string]int),
		limit: limit,
		resetTime: resetTime,
	}
	// start the reset-routine
	go rl.resetVisitorCount() // runs in the background
	return rl
}

func (rl *rateLimiter) resetVisitorCount(){
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

// Now, create the middleware( )
func (rl *rateLimiter) RateLimiterMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		rl.mu.Lock()
		defer rl.mu.Unlock()
		visitorIP:=r.RemoteAddr // minimal and simple for now.
		rl.visitors[visitorIP]++
		fmt.Printf("🔵 Visitor Count from %v is %v\n",visitorIP,rl.visitors[visitorIP])

		if rl.visitors[visitorIP] > rl.limit{
			http.Error(w, "Too Many Requests ⚠️",http.StatusTooManyRequests)
			return 
		}
		next.ServeHTTP(w,r)
	})
}
```

## 🧩 What This Middleware Does

This middleware’s job is to **limit how many requests** a single visitor (IP address) can make within a certain **time window**.

For example:

* You allow **5 requests per minute** per IP.
* If someone makes more than 5 requests before the minute resets — they get **HTTP 429 (Too Many Requests)**.

This helps **prevent abuse**, **protect your API**, and **reduce load** from bots or DoS attacks.

---

## ⚙️ Code Breakdown

### 1️⃣ Define a struct to hold rate limit info

```go
type rateLimiter struct {
	mu         sync.Mutex
	visitors   map[string]int
	limit      int
	resetTime  time.Duration
}
```

#### What each field means:

| Field       | Type             | Purpose                                                                       |
| ----------- | ---------------- | ----------------------------------------------------------------------------- |
| `mu`        | `sync.Mutex`     | Locks shared data so multiple requests don’t update the map at the same time. |
| `visitors`  | `map[string]int` | Tracks how many requests each IP has made.                                    |
| `limit`     | `int`            | Max allowed requests per IP before blocking.                                  |
| `resetTime` | `time.Duration`  | How often to reset the visitor count (e.g., every 1 minute).                  |

---

### 2️⃣ Create a new rate limiter

```go
func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	// start the reset-routine
	go rl.resetVisitorCount() 
	return rl
}
```

* This initializes a new `rateLimiter` with an empty map.
* The **goroutine** (`go rl.resetVisitorCount()`) runs a background process that periodically **clears the visitor map** — resetting request counts.

🧠 Think of it like a timer that wipes the slate clean every few minutes.

---

### 3️⃣ Background goroutine to reset visitor counts

```go
func (rl *rateLimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}
```

#### Why a **goroutine**?

* You don’t want your main server loop to block or stop while waiting for the reset interval.
* So you run this logic **concurrently** in the background.
* Every `resetTime` (e.g. 1 minute), it:

  * Sleeps for that duration.
  * Locks the map (so no one else can modify it at the same time).
  * Replaces it with a new empty map.

✅ This resets all visitor counts automatically.

---

### 4️⃣ The actual middleware

```go
func (rl *rateLimiter) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		visitorIP := r.RemoteAddr
		rl.visitors[visitorIP]++

		fmt.Printf("🔵 Visitor Count from %v is %v\n", visitorIP, rl.visitors[visitorIP])

		if rl.visitors[visitorIP] > rl.limit {
			http.Error(w, "Too Many Requests ⚠️", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
```

Here’s what happens **for every incoming request**:

| Step | Action                                 | Purpose                                                         |
| ---- | -------------------------------------- | --------------------------------------------------------------- |
| 1️⃣  | Lock the mutex (`rl.mu.Lock()`)        | Prevent multiple goroutines from writing to `visitors` at once. |
| 2️⃣  | Read the IP                            | Identify the client using `r.RemoteAddr`.                       |
| 3️⃣  | Increment the count for that IP        | Tracks how many times that client has hit the server.           |
| 4️⃣  | If over limit → return HTTP 429        | Stops further requests temporarily.                             |
| 5️⃣  | Otherwise, call `next.ServeHTTP(w, r)` | Continue processing the request normally.                       |

---

## 🧠 Why Use `sync.Mutex`

Because Go’s HTTP server is **highly concurrent** — it handles many requests at once using **goroutines**.
So, multiple requests might try to **read/write the `visitors` map** at the same time.

If two goroutines modify a map simultaneously, you’ll get a **runtime panic: concurrent map writes** ⚠️

✅ The solution: **use a Mutex (mutual exclusion lock)**.

So:

* `rl.mu.Lock()` → locks access to the map
* Code inside the lock runs safely
* `rl.mu.Unlock()` → releases the lock so others can proceed

This ensures **thread safety**.

---

## 🧵 Why Use a Goroutine

The `resetVisitorCount` method runs in an **infinite loop**, resetting the map periodically.
If you ran it in the main thread, it would **block forever** after the first `time.Sleep`.

By using:

```go
go rl.resetVisitorCount()
```

we launch it in a **background goroutine** — so it runs asynchronously while the server continues to handle requests.

---

## ⚡ Example in Action

```go
func main() {
	rl := middlewares.NewRateLimiter(5, time.Minute)

	mux := http.NewServeMux()
	mux.Handle("/", rl.RateLimiterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})))

	http.ListenAndServe(":8080", mux)
}
```

✅ Allows 5 requests/minute per IP
🚫 Returns 429 after that
♻️ Resets counts every minute

---

## 💡 TL;DR Summary

| Concept                         | Why It’s Used                                                             |
| ------------------------------- | ------------------------------------------------------------------------- |
| **`sync.Mutex`**                | To prevent concurrent writes to the shared `visitors` map.                |
| **`go rl.resetVisitorCount()`** | Runs a periodic reset in the background without blocking the main thread. |
| **`map[string]int`**            | Keeps track of each visitor’s request count by IP.                        |
| **Rate limiting**               | Prevents abuse and overuse of your API.                                   |

---

