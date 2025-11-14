package handlers

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
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