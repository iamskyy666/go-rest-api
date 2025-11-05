package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

// Track data about API (Optional middleware())
// endpoint, statuscode, duration etc.

func ResponseTimeMiddleware(next http.Handler)http.Handler{

	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Sent Response from ResponseTimeMiddleware() ✅")
		start:= time.Now()

		// create custom response-writer to capture the status-code
		wrappedWriter:= &responseWriter{ResponseWriter: w, status: http.StatusOK} // default: StatusOK

		// calculate the duration
		duration:= time.Since(start)
		w.Header().Set("X-Response-Time",duration.String())

		next.ServeHTTP(wrappedWriter,r)

		// calculate the duration
		duration= time.Since(start)

		// Log the request details
		fmt.Printf("Method: %s, URL: %s, StatusCode: %d, Duration: %v\n", r.Method,r.URL,wrappedWriter.status,duration.String())

		fmt.Println("Sent Response from Response-Writer ☑️")

	})

	
}

type responseWriter struct{
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int){
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}