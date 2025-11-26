package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/iamskyy111/go-rest-api/internal/models"
	"github.com/iamskyy111/go-rest-api/pkg/utils"
)

//! small sorting-utils f(x)
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

//! Advanced Sorting Technique (util fx)
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

//! Advanced Filtering Technique (util fx)
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

//! GET All teachers DB ops.
func GetTeachersDbHandler(teachers []models.Teacher,r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		//http.Error(w, "ERROR connecting to DATABASE ⚠️", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "ERROR connecting to DATABASE ⚠️")
	}

	qry := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
		var args []any

		// Advanced filtering f(x)
		qry, args = AddFilters(r, qry, args)

		// Advanced Sorting f(x)
		qry = AddSorting(r, qry)


	defer db.Close()


	rows, err := db.Query(qry, args...)
	if err != nil {
		//http.Error(w, "DATABASE-QUERY Error! ⚠️", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "DATABASE-QUERY Error - ERR. retrieving data! ⚠️")
	}
	defer rows.Close()

	// teacherList := make([]models.Teacher, 0)
	for rows.Next() {
		var t models.Teacher
		err := rows.Scan(&t.ID, &t.FirstName, &t.LastName, &t.Email, &t.Class, &t.Subject)
		if err != nil {
			//http.Error(w, "ERROR scanning DB-results! ⚠️", http.StatusInternalServerError)
			return nil,utils.ErrorHandler(err, "ERROR scanning DB-results! ⚠️")
		}
		teachers = append(teachers, t)
	}
	return teachers, err
}


//! GET single teacher by ID DB ops.
func GetTeacherDbHandler(id int) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		//http.Error(w, "ERROR connecting to DATABASE ⚠️", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "ERROR connecting to DATABASE ⚠️")
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

	if err == sql.ErrNoRows {
		//http.Error(w, "Teacher Not Found! ⚠️", http.StatusNotFound)
		return models.Teacher{}, utils.ErrorHandler(err,  "Teacher Not Found! ⚠️")
	} else if err != nil {
		//http.Error(w, "DB Query Error! ⚠️", http.StatusInternalServerError)
		return models.Teacher{},  utils.ErrorHandler(err,  "DB Query Error! ⚠️")
	}
	return teacher, nil
}

// Add / POST teachers DB Ops.
func AddTeachersDbHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err,  "ERROR  connecting to DATABASE ⚠️")
	}

	defer db.Close() // Don't forget to close the db.

	//stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES(?,?,?,?,?)")
	stmt, err := db.Prepare(GenerateInsertQry(models.Teacher{}))
	if err != nil {
		return nil, utils.ErrorHandler(err,  "ERROR preparing SQL Query ⚠️")
	}
	defer stmt.Close() // Don't forget to close the stmt.

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		// res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		values:= GetStructVals(newTeacher)
		res, err := stmt.Exec(values...)
		if err != nil {
			fmt.Println("ERROR:", err)
			return nil, utils.ErrorHandler(err,  "ERROR inserting DATA into DB⚠️")
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err,  "ERROR getting last-inserted-id⚠️")
		}
		newTeacher.ID = int(lastId)
		addedTeachers[i] = newTeacher
	}
	return addedTeachers, nil
}

func GenerateInsertQry(model any)string{
	modelType:=reflect.TypeOf(model)
	var columns, placeholders string
	for i := 0; i < modelType.NumField(); i++ {
		dbTag:=modelType.Field(i).Tag.Get("db")
		fmt.Println("dbTag:",dbTag)
		dbTag=strings.TrimSuffix(dbTag,",omitempty")
		if dbTag!="" && dbTag!="id"{
			// skip the ID field if it's auto-increment
			if columns!=""{
				columns+=", "
				placeholders+=", "
			}
			columns+=dbTag
			placeholders+="?"
		}
	}
	fmt.Printf("INSERT INTO teachers (%s) VALUES (%s)\n",columns,placeholders)
	return fmt.Sprintf("INSERT INTO teachers (%s) VALUES (%s)",columns,placeholders)
}

func GetStructVals(model any)[]any{
	modelVal:=reflect.ValueOf(model)
	modelType:=modelVal.Type()
	vals:=[]any{}
	for i := 0; i < modelType.NumField(); i++ {
		dbTag:=modelType.Field(i).Tag.Get("db")
		if dbTag !="" && dbTag!="id,omitempty"{
			vals= append(vals,modelVal.Field(i).Interface())
		}
	}
	log.Println("Values:",vals)
	return vals
}

//! Update/PUT teacher Db ops.
func UpdateTeachersDbHandler(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "ERROR connecting to DB ⚠️")
	}
	defer db.Close() // always close() the db.

	// extract existing info. from DB using the received id
	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id,first_name,last_name,email,class,subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	// Handle both type of errors upon Scan()
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "Teacher Not Found ⚠️")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "ERROR: Unable to retrieve data ⚠️")
	}
	updatedTeacher.ID = existingTeacher.ID

	// posting some data - Exec(), retrieving some data - Query()/QueryRow()
	_,err=db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email,updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err!= nil{
		return models.Teacher{}, utils.ErrorHandler(err, "ERROR updating teacher ⚠️",)
	}
	return updatedTeacher, nil
}


//! PATCH Multiple Teachers DB ops.
func PatchTeachersDbHandler(updates []map[string]any) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "ERROR connecting to DB ⚠️",)
	}
	defer db.Close() // always close() the db.
	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err,"ERROR starting transaction! ⚠️",)
	}

	// Access the updates
	for _, update := range updates {
		// update is a map so.. ["id"]
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			//http.Error(w, "ERROR: Invalid teacher-ID in update! ⚠️", http.StatusBadRequest)
			return utils.ErrorHandler(err,"ERROR: Invalid teacher-ID in update! ⚠️")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err,"ERROR converting id to int! ⚠️")
		}

		// instance of Teacher{}
		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return utils.ErrorHandler(err,"ERROR: Teacher not found! ⚠️")
			}
			return utils.ErrorHandler(err,"ERROR receiving teacher! ⚠️")
		}

		// Apply updates using REFLECTION
		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue // skip updating the id field
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("Cannot convert %v to %v ⚠️", val.Type(), fieldVal.Type())
							return err
						}
					}
					break
				}
			}
		}

		// Execute the update stmt
		_, err = tx.Exec(`
	UPDATE teachers 
	SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ?
	WHERE id = ?`,
			teacherFromDb.FirstName,
			teacherFromDb.LastName,
			teacherFromDb.Email,
			teacherFromDb.Class,
			teacherFromDb.Subject,
			teacherFromDb.ID,
		)

		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err,"ERROR updating teacher! ⚠️")
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err,"ERROR committing transaction! ⚠️")
	}
	return nil
}


//! Update/PUT single-teacher by ID Db ops.
func PatchSingleTeacherDbOps(id int, updates map[string]any) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err,"ERROR connecting to DB ⚠️")
	}
	defer db.Close() // always close() the db.

	// execute the query to find the teacher
	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id,first_name,last_name,email,class,subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	// Handle both type of errors upon Scan()
	if err == sql.ErrNoRows {
		return models.Teacher{},utils.ErrorHandler(err,"Teacher Not Found ⚠️")
	} else if err != nil {
		return models.Teacher{

			//! 💡 apply updates - refactored, using reflect pkg.
		}, utils.ErrorHandler(err,"Unable to retrieve data ⚠️")
	}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			field.Tag.Get("json")
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					fieldVal := teacherVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	// send existingTeacher{} back to the DB for updation
	// posting some data - Exec(), retrieving some data - Query()/QueryRow()
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err,"ERROR updating teacher ⚠️")
	}
	return existingTeacher, nil
}


//! Delete Single Teacher Db ops.
func DeleteSingleTeacherDbHandler( id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err,"ERROR connecting to DB ⚠️")
	}
	defer db.Close() // always close() the db.

	// res/result - confirmation
	res, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err,"ERROR deleting teacher ⚠️")
	}

	fmt.Println(res.RowsAffected()) // for our info.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err,"ERROR deleting teacher ⚠️")
	}

	if rowsAffected == 0 {
		return utils.ErrorHandler(err,"ERROR retrieving deleted-teacher ⚠️")
	}
	return nil
}

//! Delete Multiple Teachers Db ops.
func DeleteTeachersDbHandler(ids []int) ([]int, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err,"ERROR connecting to DB ⚠️")
	}
	defer db.Close() // always close() the db.

	// Transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err,"ERROR starting transaction ⚠️")
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err,"ERROR preparing DELETE statement ⚠️")
	}

	defer stmt.Close() // Always close stmt

	deletedIds := []int{}

	for _, id := range ids {
		res, err := stmt.Exec(int64(id))
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err,"ERROR deleting teachers! ⚠️")
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, utils.ErrorHandler(err,"ERROR retrieving deleted-teachers ⚠️")
		}

		// if teacher was deleted, then add the teacher-ID to the deletedIDs slice[]
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1 {
			tx.Rollback()
			return nil, utils.ErrorHandler(err,fmt.Sprintf("ID %d does not exist ⚠️", id))
		}
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err,"ERROR committing transaction! ⚠️")
	}

	if len(deletedIds) < 1 {
		return nil, utils.ErrorHandler(err,"IDs do not exist ⚠️")
	}
	return deletedIds, nil
}