package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/iamskyy111/go-rest-api/internal/api/middlewares"
)

func RootHandler(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Root-Route ✅")
		w.Write([]byte("Hello Root-Route ✅"))
		fmt.Println("Root Route ✅", r.Method)
	}

func TeachersHandler(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Teachers-Route ✅")

		switch r.Method {
				case http.MethodGet:
			w.Write([]byte("Hello GET method on Teachers-Route ✅"))
			fmt.Println("Hello GET method on Teachers-Route ✅")
			return
				case http.MethodPost:							
			w.Write([]byte("Hello Post method on Teachers-Route ✅"))
			fmt.Println("Hello Post method on Teachers-Route ✅")
			return
				case http.MethodPatch:
			w.Write([]byte("Hello Patch method on Teachers-Route ✅"))
			fmt.Println("Hello Patch method on Teachers-Route ✅")
			return
				case http.MethodDelete:
			w.Write([]byte("Hello Delete method on Teachers-Route ✅"))
			fmt.Println("Hello Delete method on Teachers-Route ✅")
			return	
				default:
			w.Write([]byte("Hello UNKNOWN method on Teachers-Route!"))
			fmt.Println("Hello UNKNOWN method on Teachers-Route !")
			return	

		}

	}	

func StudentsHandler(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello Students-Route ✅"))
		fmt.Println("Students Route ✅", r.Method)
	}	

func ExecsHandler(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello Execs-Route ✅"))
		fmt.Println("Execs Route ✅", r.Method)
	}

func main() {
	PORT := ":3000"
	cert:= "cert.pem"
	key:="key.pem"

	mux:= http.NewServeMux()

	mux.HandleFunc("/", RootHandler )
	mux.HandleFunc("/teachers/", TeachersHandler)
	mux.HandleFunc("/students", StudentsHandler)
	mux.HandleFunc("/execs", ExecsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// initialize the rate-limiter
	rl:= middlewares.NewRateLimiter(5, time.Minute)

	// instance of the HppOptions struct
	hppOptions:= middlewares.HPPOptions{
		CheckQuery: true,
		CheckBody: true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		WhiteList: []string{"sortBy","sortOrder","name","age","class"},
	}

	// Efficient Middleware Ordering/Chaining ✅
	secureMux:= middlewares.CorsMiddleware(rl.RateLimiterMiddleware(middlewares.ResponseTimeMiddleware(middlewares.SecurityHeaders(middlewares.CompressionMiddleware(middlewares.HppMiddleware(hppOptions)(mux))))))



	// Create custom-server
	server:= &http.Server{
		Addr:PORT,
		Handler: secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on PORT", PORT,"🟢")
	err:= server.ListenAndServeTLS(cert,key)
	if err!=nil{
		log.Fatal("⚠️ERROR. starting the server:",err)
	}
}
