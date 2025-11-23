package main

import (
	"fmt"
	"log"
	"net/http"
)

//💡 modern way of writing routers

func main() {
	mux:= http.NewServeMux()

	// method-based routing
	mux.HandleFunc("POST /items/create", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w,"Item Created ✅")
	})

	mux.HandleFunc("DELETE /items/delete", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w,"Item Deleted ☑️")
	})

	// wildcard in pattern - path parameter
	mux.HandleFunc("GET /teachers/{id}", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w,"🔵 Teacher ID: %s",r.PathValue("id"))
	})

	// wildcard with  "..." pattern
	mux.HandleFunc("/files/{path...}", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w,"🟡 Path: %s",r.PathValue("path"))
	})

	// mix n match
	mux.HandleFunc("/path1/{param1}", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w,"🟠 Param1: %s",r.PathValue("param1"))
	})

	// BEWARE ❌❌❌
	// mux.HandleFunc("/param2/{path2}", func(w http.ResponseWriter, r *http.Request){
	// 	fmt.Fprintf(w,"🟢 Param2: %s",r.PathValue("param2"))
	// })

	mux.HandleFunc("/path1/path2", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w,"🟢 Param2: %s",r.PathValue("param2"))
	})


	err:=http.ListenAndServe(":8080",mux)
	if err != nil {
		fmt.Println("ERROR:",err)
		log.Fatal(err)
	}
}