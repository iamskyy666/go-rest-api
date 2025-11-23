package router

import (
	"net/http"

	"github.com/iamskyy111/go-rest-api/internal/api/handlers"
)

func Router() *http.ServeMux{
mux:= http.NewServeMux()

mux.HandleFunc("/", handlers.RootHandler )

// Teachers Handlers()
mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
mux.HandleFunc("POST /teachers/", handlers.AddTeacherHandler)
mux.HandleFunc("PUT /teachers/", handlers.UpdateTeacherHandler)
mux.HandleFunc("PATCH /teachers/", handlers.PatchTeacherHandler)
mux.HandleFunc("DELETE /teachers/", handlers.DeleteTeacherHandler)

mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandler)
mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler)
mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

mux.HandleFunc("/students", handlers.StudentsHandler)
mux.HandleFunc("/execs", handlers.ExecsHandler)

return mux
}

