package router

import (
	"net/http"

	"github.com/iamskyy111/go-rest-api/internal/api/handlers"
)

func Router() *http.ServeMux{
mux:= http.NewServeMux()

mux.HandleFunc("/", handlers.RootHandler )
mux.HandleFunc("/teachers/", handlers.TeachersHandler)
mux.HandleFunc("/students", handlers.StudentsHandler)
mux.HandleFunc("/execs", handlers.ExecsHandler)

return mux
}

