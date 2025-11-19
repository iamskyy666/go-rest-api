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

//! 1️⃣☑️ GET/FETCH teacher(s)
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "ERROR connecting to DATABASE ⚠️", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	// If no ID: GET ALL TEACHERS
	if idStr == "" {

		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		qry := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
		var args []any

		if firstName != "" {
			qry += " AND first_name = ?"
			args = append(args, firstName)
		}

		if lastName != "" {
			qry += " AND last_name = ?"
			args = append(args, lastName)
		}

		rows, err := db.Query(qry, args...)
		if err != nil {
			http.Error(w, "DATABASE-QUERY Error! ⚠️", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		teacherList := make([]models.Teacher, 0)
		for rows.Next() {
			var t models.Teacher
			err := rows.Scan(&t.ID, &t.FirstName, &t.LastName, &t.Email, &t.Class, &t.Subject)
			if err != nil {
				http.Error(w, "ERROR scanning DB-results! ⚠️", http.StatusInternalServerError)
				return
			}
			teacherList = append(teacherList, t)
		}

		resp := struct {
			Status string            `json:"status"`
			Count  int               `json:"count"`
			Data   []models.Teacher  `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return   // ← CRUCIAL FIX
	}

	// SINGLE TEACHER ======================================
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

	if err == sql.ErrNoRows {
		http.Error(w, "Teacher Not Found! ⚠️", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "DB Query Error! ⚠️", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

//! 2️⃣☑️ ADD/POST Teacher(s)
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



// Handler
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