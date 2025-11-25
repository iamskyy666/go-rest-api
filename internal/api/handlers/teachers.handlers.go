package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/iamskyy111/go-rest-api/internal/models"
	"github.com/iamskyy111/go-rest-api/internal/repositories/sqlconnect"
)

// 💡 All Ops. apart from GET requires db.Exec()

// CRUD ⭐
//! 1️⃣☑️ GET/FETCH teacher(s)
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
		var teachers []models.Teacher
		teachers, err := sqlconnect.GetTeachersDbHandler(teachers,r) // db ops.
		if err!=nil{
			return
		}

		resp := struct {
			Status string            `json:"status"`
			Count  int               `json:"count"`
			Data   []models.Teacher  `json:"data"`
		}{
			Status: "success",
			Count:  len(teachers),
			Data:   teachers,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
}


//! 2️⃣☑️ GET/FETCH single-teacher /id
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	// SINGLE TEACHER ======================================
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	teacher, err := sqlconnect.GetTeacherDbHandler(id)
	if err!=nil {
		fmt.Println("ERROR:",err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}


//! 3️⃣☑️ ADD/POST Teacher(s)
func AddTeacherHandler(w http.ResponseWriter, r *http.Request){

	var newTeachers []models.Teacher
	err:=json.NewDecoder(r.Body).Decode(&newTeachers) // we can add 1 or multiple values in a list
	if err != nil {
		http.Error(w,"Invalid Request Body!",http.StatusBadRequest)
		return
	}

	// Connect to DB
	addedTeachers, err := sqlconnect.AddTeachersDbHandler(newTeachers)
	if err!=nil {
		fmt.Println("ERROR:",err)
		return
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
	idStr:= r.PathValue("id")
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
	updatedTeacherFromDb, err := sqlconnect.UpdateTeachersDbHandler(id,updatedTeacher)
	if err!=nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(updatedTeacherFromDb)
}

//! 4️⃣☑️ Partially-Edit/PATCH Single Teacher/id
func PatchTeacherHandler(w http.ResponseWriter, r *http.Request){
	// get id from the params and convert it to an 'int'
	idStr:= r.PathValue("id")
	id,err:= strconv.Atoi(idStr)
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"Invalid teacher-ID ⚠️",http.StatusBadRequest)
		return
	}

	// decode JSON body into the model
	var updates map[string]any
	err=json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"Invalid request-payload ⚠️",http.StatusBadRequest)
		return
	}

	// connect to the DB
	updatedteacher, err := sqlconnect.PatchSingleTeacherDbOps(id, updates)
	if err!=nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(updatedteacher)
}

//! 5️⃣☑️ PATCH Multiple-Teachers
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request){
	var updates []map[string]any
	err:=json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"ERROR: Invalid request-payload ⚠️",http.StatusBadRequest)
		return
	}

	// connect to the DB
	err = sqlconnect.PatchTeachersDbHandler(updates)
	if err!=nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent) 
}


//! 6️⃣☑️ DELETE Single Teacher/id
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request){
	// extract id from the params and convert it to an 'int'
	idStr:= r.PathValue("id")
	id,err:= strconv.Atoi(idStr)
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"Invalid teacher-ID ⚠️",http.StatusBadRequest)
		return
	}

	// connect to the DB
	err = sqlconnect.DeleteSingleTeacherDbHandler(id)
	if err!=nil {
		log.Println(err)
		return
	}

	// no-response if delete is successful
	// w.WriteHeader(http.StatusNoContent)

	// Alternatively, response-body (better approach!)
	w.Header().Set("Content-Type","application/json")
	response:=struct{
		Status string `json:"status"`
		ID int `json:"id"`
	}{
		Status: "Teacher successfully DELETED ✅",
		ID: id,
	}

	json.NewEncoder(w).Encode(response)
}


 //! 7️⃣☑️ DELETE Multiple-Teachers
 func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request){
	var ids []int
	err:=json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println("ERROR:",err)
		http.Error(w,"ERROR: Invalid request-payload ⚠️",http.StatusBadRequest)
		return
	}

	// connect to the DB
	deletedIds, err := sqlconnect.DeleteTeachersDbHandler(ids)
	if err!=nil {
		log.Println("ERROR:",err)
		return
	}

	// Finally... Send the response
	w.Header().Set("Content-Type","application/json")
	resp:=struct{
		Status string `json:"status"`
		DeletedIDs []int `json:"deleted_ids"`
	}{
		Status: "Teachers Successfully Deleted ✅",
		DeletedIDs: deletedIds,
	}

	json.NewEncoder(w).Encode(resp)
}

