package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/iamskyy111/go-rest-api/internal/models"
	"github.com/iamskyy111/go-rest-api/internal/repositories/sqlconnect"
)

// small sorting-utils f(x)
func IsValidSortOrder(order string)bool{
	return order=="asc" || order=="desc"
}

func IsValidSortField(field string)bool{
	//💡 Map[][] is faster than a []slice or an []array
	validFields:=map[string]bool{
		"first_name":true,
		"last_name":true,
		"email":true,
		"class":true,
		"subject":true,
	}
	return validFields[field]
}


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
		qry := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
		var args []any

		// firstName := r.URL.Query().Get("first_name")
		// lastName := r.URL.Query().Get("last_name")

		// Advanced filtering f(x)
		qry, args = AddFilters(r, qry, args)

		// Advanced Sorting f(x)
		qry = AddSorting(r, qry)

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

// Advanced Sorting Technique
func AddSorting(r *http.Request, qry string) string {
	sortParams := r.URL.Query()["sortby"]

	// if params are not empty, then..
	if len(sortParams) > 0 {
		qry += " ORDER BY"

		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !IsValidSortField(field) || !IsValidSortOrder(order) {
				continue
			}
			if i > 0 {
				qry += ","
			}
			qry += " " + field + " " + order
		}
	}
	return qry
}

// Advanced Filtering Technique
func AddFilters(r *http.Request, qry string, args []any) (string, []any) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		val := r.URL.Query().Get(param)
		if val != "" {
			qry += " AND " + dbField + " = ?"
			args = append(args, val)
		}
	}
	return qry, args
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


//! 3️⃣☑️ UPDATE/PUT Teachers/id
//💡 PUT replaces the whole entry, leaving 1 blank will also result in a blank-entry in the db (unlike PATCH/partial update)
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request){
	// get id from the params and convert it to an 'int'
	idStr:= strings.TrimPrefix(r.URL.Path,"/teachers/")
	id,err:= strconv.Atoi(idStr)
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"Invalid teacher-ID ⚠️",http.StatusBadRequest)
		return
	}

	// decode JSON body into the model
	var updatedTeacher models.Teacher
	err=json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"Invalid request-payload ⚠️",http.StatusBadRequest)
		return
	}

	// connect to the DB
	db,err:=sqlconnect.ConnectDB()
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"ERROR connecting to DB ⚠️",http.StatusInternalServerError)
		return
	}
	defer db.Close() // always close() the db.

	// extract existing info. from DB using the received id
	var existingTeacher models.Teacher
	err=db.QueryRow("SELECT id,first_name,last_name,email,class,subject FROM teachers WHERE id = ?",id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	// Handle both type of errors upon Scan()
	if err== sql.ErrNoRows{
		log.Println("ERROR:",err)
		http.Error(w,"Teacher Not Found ⚠️",http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"Unable to retrieve data ⚠️",http.StatusInternalServerError)
		return
	}

	updatedTeacher.ID = existingTeacher.ID

	// posting some data - Exec(), retrieving some data - Query()/QueryRow()
	_,err=db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email,updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err!= nil{
		log.Println("ERROR:",err)
		http.Error(w,"ERROR updating teacher ⚠️",http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}


//! 4️⃣☑️ Partially-Edit/PATCH Teachers/id
func PatchTeacherHandler(w http.ResponseWriter, r *http.Request){
	// get id from the params and convert it to an 'int'
	idStr:= strings.TrimPrefix(r.URL.Path,"/teachers/")
	id,err:= strconv.Atoi(idStr)
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"Invalid teacher-ID ⚠️",http.StatusBadRequest)
		return
	}

	// decode JSON body into the model
	var updates map[string]interface{}
	err=json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"Invalid request-payload ⚠️",http.StatusBadRequest)
		return
	}

	// connect to the DB
	db,err:=sqlconnect.ConnectDB()
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"ERROR connecting to DB ⚠️",http.StatusInternalServerError)
		return
	}
	defer db.Close() // always close() the db.

	// execute the query to find the teacher
	var existingTeacher models.Teacher
	err=db.QueryRow("SELECT id,first_name,last_name,email,class,subject FROM teachers WHERE id = ?",id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	// Handle both type of errors upon Scan()
	if err== sql.ErrNoRows{
		log.Println("ERROR:",err)
		http.Error(w,"Teacher Not Found ⚠️",http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"Unable to retrieve data ⚠️",http.StatusInternalServerError)
		return
	}

	// apply updates
	// for k,v:=range updates{
	// 	switch k {
	// 		case "first_name":
	// 		existingTeacher.FirstName = v.(string)
	// 		case "last_name":
	// 		existingTeacher.LastName = v.(string)	
	// 		case "email":
	// 		existingTeacher.Email = v.(string)	
	// 		case "class":
	// 		existingTeacher.Class = v.(string)	
	// 		case "subject":
	// 		existingTeacher.Subject = v.(string)
	// 	}
	// }

	//! 💡 apply updates - refactored, using reflect pkg.
	teacherVal:= reflect.ValueOf(&existingTeacher).Elem()
	teacherType:= teacherVal.Type()

	for k,v:=range updates{
		for i:=0; i<teacherVal.NumField();i++{
			field:=teacherType.Field(i)
			field.Tag.Get("json")
			if field.Tag.Get("json")==k+",omitempty" {
				if teacherVal.Field(i).CanSet(){
					fieldVal:=teacherVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	// send existingTeacher{} back to the DB for updation
	// posting some data - Exec(), retrieving some data - Query()/QueryRow()
	_,err=db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email,existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err!= nil{
		log.Println("ERROR:",err)
		http.Error(w,"ERROR updating teacher ⚠️",http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(existingTeacher)
}



// Handlers
func TeachersHandler(w http.ResponseWriter, r *http.Request){
		switch r.Method {
			case http.MethodGet:
			GetTeacherHandler(w,r)	
			return
			case http.MethodPost:							
			AddTeacherHandler(w,r)
			return
			case http.MethodPut:							
			UpdateTeacherHandler(w,r)
			return
			case http.MethodPatch:
			PatchTeacherHandler(w,r)
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