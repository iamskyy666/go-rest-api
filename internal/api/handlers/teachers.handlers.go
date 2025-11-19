package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/iamskyy111/go-rest-api/internal/models"
	"github.com/iamskyy111/go-rest-api/internal/repositories/sqlconnect"
)

var teachers = make(map[int]models.Teacher)
// var mutex = &sync.Mutex{}
// var nextID = 1

// Initialize some/dummy data 🔵
// func init(){
// 		teachers[nextID] = models.Teacher{
// 		ID: nextID,
// 		FirstName: "John",
// 		LastName: "Doe",
// 		Class: "9A",
// 		Subject: "Math",
// 	}
// 	nextID++
// 		teachers[nextID] = models.Teacher{
// 		ID: nextID,
// 		FirstName: "Jane",
// 		LastName: "Smith",
// 		Class: "10A",
// 		Subject: "English Lit.",
// 	}
// 	nextID++
// 		teachers[nextID] = models.Teacher{
// 		ID: nextID,
// 		FirstName: "Jane",
// 		LastName: "Doe",
// 		Class: "11A",
// 		Subject: "Geography",
// 	}
// 	nextID++
// }

// ☑️ GET single teacher based on ID
func GetTeacherHandler(w http.ResponseWriter, r *http.Request){
	// Connect to DB
	db,err:=sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "ERROR  connecting to DATABASE ⚠️",http.StatusInternalServerError)
		return
	}

	defer db.Close() // Don't forget to close the db.

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
	// teacher,exists := teachers[id]
	// if !exists{
	// 	http.Error(w, "⚠️Teacher Not Found!",http.StatusNotFound)
	// 	return
	// }
	var teacher models.Teacher
	err=db.QueryRow("SELECT id, first_name, last_name,email, class, subject FROM teachers WHERE id=?",id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher Not Found! ⚠️",http.StatusNotFound)
		return
	}else if err!=nil{
		http.Error(w, "DB Query Error! ⚠️",http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(teacher)
}

// ☑️ ADD/POST Teacher(s)
func AddTeacherHandler(w http.ResponseWriter, r *http.Request){
	// Connect to DB
	db,err:=sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "ERROR  connecting to DATABASE ⚠️",http.StatusInternalServerError)
		return
	}

	defer db.Close() // Don't forget to close the db.

	var newTeachers []models.Teacher
	err=json.NewDecoder(r.Body).Decode(&newTeachers) // we can add 1 or multiple values in a list
	if err != nil {
		http.Error(w,"Invalid Request Body!",http.StatusBadRequest)
		return
	}

	stmt,err:=db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES(?,?,?,?,?)")
	if err != nil {
		http.Error(w,"ERROR preparing SQL Query ⚠️",http.StatusInternalServerError)
		return
	}
	defer stmt.Close() // Don't forget to close the stmt.

	addedTeachers:=make([]models.Teacher,len(newTeachers))
	for i, newTeacher:= range newTeachers{
		res,err:=stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
		http.Error(w,"ERROR inserting DATA into DB",http.StatusInternalServerError)
		fmt.Println("ERROR:",err)
		return
	}
	lastId,err := res.LastInsertId()
		if err != nil {
		http.Error(w,"ERROR getting last-inserted-id⚠️",http.StatusInternalServerError)
		return
	}
	newTeacher.ID = int(lastId)
	addedTeachers[i] = newTeacher
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
		GetTeacherHandler(w,r)	
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