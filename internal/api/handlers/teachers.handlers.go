package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/iamskyy111/go-rest-api/internal/models"
)

var teachers = make(map[int]models.Teacher)
var mutex = &sync.Mutex{}
var nextID = 1

// Initialize some/dummy data
func init(){
		teachers[nextID] = models.Teacher{
		ID: nextID,
		FirstName: "John",
		LastName: "Doe",
		Class: "9A",
		Subject: "Math",
	}
	nextID++
		teachers[nextID] = models.Teacher{
		ID: nextID,
		FirstName: "Jane",
		LastName: "Smith",
		Class: "10A",
		Subject: "English Lit.",
	}
	nextID++
		teachers[nextID] = models.Teacher{
		ID: nextID,
		FirstName: "Jane",
		LastName: "Doe",
		Class: "11A",
		Subject: "Geography",
	}
	nextID++
}

// ☑️ GET all teachers
func GetTeachersHandler(w http.ResponseWriter, r *http.Request){
	path := strings.TrimPrefix(r.URL.Path,"/teachers/")
	idStr:=strings.TrimSuffix(path,"/")

	fmt.Println(idStr)

	if idStr==""{
	firstName:= r.URL.Query().Get("first_name")
	lastName:= r.URL.Query().Get("last_name")
	teacherList:= make([]models.Teacher,0,len(teachers))
	for _, teacher:= range teachers{
		if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName){
			teacherList = append(teacherList, teacher)
		}
	}
	

	resp:= struct{
		Status string `json:"status"`
		Count int `json:"count"`
		Data []models.Teacher `json:"data"`
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

// ☑️ ADD/POST Teacher(s)
func AddTeacherHandler(w http.ResponseWriter, r *http.Request){
	// Use mutex when adding a new teacher / POST
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []models.Teacher
	err:=json.NewDecoder(r.Body).Decode(&newTeachers) // we can add 1 or multiple values in a list
	if err != nil {
		http.Error(w,"Invalid Request Body!",http.StatusBadRequest)
		return
	}

	addedTeachers:=make([]models.Teacher,len(newTeachers))
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
		Data []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count: len(addedTeachers),
		Data: addedTeachers,
	}
	json.NewEncoder(w).Encode(resp)
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