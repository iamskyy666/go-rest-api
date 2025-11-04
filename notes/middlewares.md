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
