package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/iamskyy111/go-rest-api/internal/api/middlewares"
)

// 💡 models
type Teacher struct{
	ID int
	FirstName string
	LastName string
	Class string
	Subject string
}

var teachers = make(map[int]Teacher)
var mutex= &sync.Mutex{}
var nextID = 1

// Initialize some data
func init(){
	teachers[nextID] = Teacher{
		ID: nextID,
		FirstName: "John",
		LastName: "Doe",
		Class: "9A",
		Subject: "Math",
	}
	nextID++
	teachers[nextID] = Teacher{
			ID: nextID,
		FirstName: "Jane",
		LastName: "Smith",
		Class: "10A",
		Subject: "English Lit.",
	}
	nextID++
		teachers[nextID] = Teacher{
		ID: nextID,
		FirstName: "Jane",
		LastName: "Doe",
		Class: "11A",
		Subject: "Geography",
	}
}

// Initially, in-memory handler functions.
// ☑️ GET all teachers
func GetTeachersHandler(w http.ResponseWriter, r *http.Request){


	path := strings.TrimPrefix(r.URL.Path,"/teachers/")
	idStr:=strings.TrimSuffix(path,"/")

	fmt.Println(idStr)

	if idStr==""{
	firstName:= r.URL.Query().Get("first_name")
	lastName:= r.URL.Query().Get("lst_name")
	teacherList:= make([]Teacher,0,len(teachers))
	for _, teacher:= range teachers{
		if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName){
			teacherList = append(teacherList, teacher)
		}
	}
	

	resp:= struct{
		Status string `json:"status"`
		Count int `json:"count"`
		Data []Teacher `json:"data"`
	}{
		Status: "success",
		Count: len(teacherList),
		Data: teacherList,
	}

	// send the resp. in JSON{}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(resp)
}
	// Handle Path Param
	id,err:= strconv.Atoi(idStr)
	if err!=nil{
		fmt.Println("⚠️ERROR:",err)
		return
	}
	teacher,exists := teachers[id]
	if !exists{
		http.Error(w, "⚠️Teacher Not Found!",http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(teacher)
}

func RootHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
				case http.MethodGet:
			w.Write([]byte("Hello GET method on Root-Route ✅"))
			fmt.Println("Hello GET method on Root-Route ✅")
			return
				case http.MethodPost:							
			w.Write([]byte("Hello Post method on Root-Route ✅"))
			fmt.Println("Hello Post method on Root-Route ✅")
			return
				case http.MethodPatch:
			w.Write([]byte("Hello Patch method on Root-Route ✅"))
			fmt.Println("Hello Patch method on Root-Route ✅")
			return
				case http.MethodDelete:
			w.Write([]byte("Hello Delete method on Root-Route ✅"))
			fmt.Println("Hello Delete method on Root-Route ✅")
			return	
				default:
			w.Write([]byte("Hello UNKNOWN method on Root-Route!"))
			fmt.Println("Hello UNKNOWN method on Root-Route !")
			return	
		}
	}

func TeachersHandler(w http.ResponseWriter, r *http.Request){
		switch r.Method {
				case http.MethodGet:
		GetTeachersHandler(w,r)	
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
		switch r.Method {
				case http.MethodGet:
			w.Write([]byte("Hello GET method on Students-Route ✅"))
			fmt.Println("Hello GET method on Students-Route ✅")
			return
				case http.MethodPost:							
			w.Write([]byte("Hello Post method on Students-Route ✅"))
			fmt.Println("Hello Post method on Students-Route ✅")
			return
				case http.MethodPatch:
			w.Write([]byte("Hello Patch method on Students-Route ✅"))
			fmt.Println("Hello Patch method on Students-Route ✅")
			return
				case http.MethodDelete:
			w.Write([]byte("Hello Delete method on Students-Route ✅"))
			fmt.Println("Hello Delete method on Students-Route ✅")
			return	
				default:
			w.Write([]byte("Hello UNKNOWN method on Students-Route!"))
			fmt.Println("Hello UNKNOWN method on Students-Route !")
			return	
		}
	}	

func ExecsHandler(w http.ResponseWriter, r *http.Request){
		switch r.Method {
				case http.MethodGet:
			w.Write([]byte("Hello GET method on Execs-Route ✅"))
			fmt.Println("Hello GET method on Execs-Route ✅")
			return
				case http.MethodPost:							
			w.Write([]byte("Hello Post method on Execs-Route ✅"))
			fmt.Println("Hello Post method on Execs-Route ✅")
			return
				case http.MethodPatch:
			w.Write([]byte("Hello Patch method on Execs-Route ✅"))
			fmt.Println("Hello Patch method on Execs-Route ✅")
			return
				case http.MethodDelete:
			w.Write([]byte("Hello Delete method on Execs-Route ✅"))
			fmt.Println("Hello Delete method on Execs-Route ✅")
			return	
				default:
			w.Write([]byte("Hello UNKNOWN method on Execs-Route!"))
			fmt.Println("Hello UNKNOWN method on Execs-Route !")
			return	
		}
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

	// initialize the rate-limiter ✅
	//rl:= middlewares.NewRateLimiter(5, time.Minute)

	// instance of the HppOptions struct ✅
	// hppOptions:= middlewares.HPPOptions{
	// 	CheckQuery: true,
	// 	CheckBody: true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	WhiteList: []string{"sortBy","sortOrder","name","age","class"},
	// }

	// secureMux:= middlewares.CorsMiddleware(rl.RateLimiterMiddleware(middlewares.ResponseTimeMiddleware(middlewares.SecurityHeaders(middlewares.CompressionMiddleware(middlewares.HppMiddleware(hppOptions)(mux))))))

	// Efficient Middleware Ordering/Chaining ✅
	// secureMux:= applyMiddlewares(mux, middlewares.HppMiddleware(hppOptions), middlewares.CompressionMiddleware, middlewares.CompressionMiddleware, middlewares.SecurityHeaders, middlewares.ResponseTimeMiddleware, rl.RateLimiterMiddleware,middlewares.CorsMiddleware)

	secureMux:= middlewares.SecurityHeaders(mux)



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

// Efficient Middleware Chaining 💡
// Middleware is a function that wraps http.handler with additional functionalities

type Middleware func(http.Handler)http.Handler
func ApplyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler{
	for _,mw:= range middlewares{
		handler = mw(handler)
	}
	return handler
}