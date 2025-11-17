package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/iamskyy111/go-rest-api/internal/api/middlewares"
	"github.com/iamskyy111/go-rest-api/internal/api/router"
	"github.com/iamskyy111/go-rest-api/internal/repositories/sqlconnect"
	"github.com/joho/godotenv"
)


func main() {
	// godotenv loads the .env file vars as if they're part of the system OS.
	// mentioning it in the main() is enough
	err:=godotenv.Load()
	if err != nil {
		fmt.Println("⚠️Error loading .env files:",err)
		return
	}

	// DB connection
	_,err=sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("ERROR:",err)
	}


	PORT := os.Getenv("API_PORT")
	cert:= "cert.pem"
	key:="key.pem"



	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	router:=router.Router()
	secureMux:= middlewares.SecurityHeaders(router)

	// Create custom-server
	server:= &http.Server{
		Addr:PORT,
		Handler: secureMux,
		TLSConfig: tlsConfig,
	}
	fmt.Println("Server is running on PORT", PORT,"🟢")
	err= server.ListenAndServeTLS(cert,key)
	if err!=nil{
		log.Fatal("⚠️ERROR. starting the server:",err)
	}
}


// initialize the rate-limiter ✅
	//rl:= middlewares.NewRateLimiter(5, time.Minute)

	// instance of the HppOptions struct ✅
	// hppOptions:= middlewares.HPPOptions{
	// 	CheckQuery: true,
	// 	CheckBody: true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	WhiteList: []string{"sortBy","sortOrder","name","age","class"},
	// }

	//secureMux:= utils.ApplyMiddlewares(mux, middlewares.HppMiddleware(hppOptions), middlewares.CompressionMiddleware, middlewares.CompressionMiddleware, middlewares.SecurityHeaders, middlewares.ResponseTimeMiddleware, rl.RateLimiterMiddleware,middlewares.CorsMiddleware)
