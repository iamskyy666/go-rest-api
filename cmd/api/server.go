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
	ID int `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Class string `json:"class,omitempty"`
	Subject string `json:"subject,omitempty"`
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
	nextID++
}

// Initially, in-memory handler functions.
// ☑️ GET all teachers
func GetTeachersHandler(w http.ResponseWriter, r *http.Request){


	path := strings.TrimPrefix(r.URL.Path,"/teachers/")
	idStr:=strings.TrimSuffix(path,"/")

	fmt.Println(idStr)

	if idStr==""{
	firstName:= r.URL.Query().Get("first_name")
	lastName:= r.URL.Query().Get("last_name")
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

// ☑️ ADD Teacher(s)
func AddTeacherHandler(w http.ResponseWriter, r *http.Request){
	// Use mutex when adding a new teacher / POST
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []Teacher
	err:=json.NewDecoder(r.Body).Decode(&newTeachers) // we can add 1 or multiple values in a list
	if err != nil {
		http.Error(w,"Invalid Request Body!",http.StatusBadRequest)
		return
	}

	addedTeachers:=make([]Teacher,len(newTeachers))
	for i, newTeacher:= range newTeachers{
		newTeacher.ID = nextID
		teachers[nextID] = newTeacher
		addedTeachers[i] = newTeacher
		nextID++
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	resp:= struct{
		Status string `json:"status"`
		Count int `json:"count"`
		Data []Teacher `json:"data"`
	}{
		Status: "success",
		Count: len(addedTeachers),
		Data: addedTeachers,
	}
	json.NewEncoder(w).Encode(resp)
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
			AddTeacherHandler(w,r)
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

type Middleware func(http.Handler)http.Handler
func ApplyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler{
	for _,mw:= range middlewares{
		handler = mw(handler)
	}
	return handler
}