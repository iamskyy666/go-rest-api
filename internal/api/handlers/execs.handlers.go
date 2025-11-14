package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
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