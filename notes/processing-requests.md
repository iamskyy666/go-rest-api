
---
# PROCESSING REQUESTS: LEGACY WAY ⌛🛜
```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// For multiple usecases, use structs (MODELs) instead of simple-maps
type Country struct{
	Name string `json:"country"`
	Capital string `json:"capital"`
	Language string `json:"language"`
	ISD_CODE string `json:"isd-code"`
}

func main() {
	PORT := ":3000"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Root-Route ✅")
		w.Write([]byte("Hello Root-Route ✅"))
		fmt.Println("Root Route ✅", r.Method)
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Teachers-Route ✅")

		switch r.Method {
				case http.MethodGet:
			w.Write([]byte("Hello GET method on Teachers-Route ✅"))
			fmt.Println("Hello GET method on Teachers-Route ✅")
			return
				case http.MethodPost:
			//💡 Parse form data (necessary for x-www-form-urlencoded)
			
			err:=r.ParseForm()
			if err!=nil{
				http.Error(w,"ERROR parsing form!",http.StatusBadRequest)
				return
			}

			fmt.Println("FORM:",r.Form) // FORM: map[age:[30] city:[CCU, BER] name:[Skyy]]

			// Prepare the response-data
			resp:= make(map[string]any)
			for k,v:= range r.Form{
				// resp[k]=v  // Processed Response Map: map[age:[30] city:[CCU,BER,] name:[Skyy]]
				resp[k]=v[0]  // Processed Response Map: map[age:30 city:CCU,BER, name:Skyy]
			}
			fmt.Println("Processed Response Map:",resp)


			//💡 RAW Body { }
			body,err:=io.ReadAll(r.Body)
			if err!=nil{
				return
			}
			r.Body.Close() //⚠️ remember to close the body, bcz, field Body io.ReadCloser!
			fmt.Println("RAW Body:",body)
			/*
			RAW Body: [123 13 10 32 32 32 32 34 99 111 117 110 116 114 121 34 58 32 34 76 117 120 101 109 98 117 114 103 34 44 13 10 32 32 32 32 34 99 97 112 105 116 97 108 34 58 32 34 76 117 120 101 109 98 117 114 103 32 67 105 116 121 34 44 13 10 32 32 32 32 34 108 97 110 103 117 97 103 101 34 58 32 34 76 195 171 116 122 101 98 117 101 114 103 101 115 99 104 34 44 13 10 32 32 32 32 34 105 115 100 45 99 111 100 101 34 58 34 43 51 53 50 34 13 10 125]
			*/

			fmt.Println("Processed Body:", string(body))
			/*
			Processed Body: {
				"country": "Luxemburg",
				"capital": "Luxemburg City",
				"language": "Lëtzebuergesch",
				"isd-code":"+352"
				}
			*/

			// 💡 UNMARSHAL RAW-Body in JSON format:
			var Lëtzebuerg Country
			err=json.Unmarshal(body, &Lëtzebuerg)
			if err!=nil{
				return
			}
			fmt.Println("Lëtzebuerg:",Lëtzebuerg) // Lëtzebuerg: {Luxembourg Luxembourg-City Lëtzebuergesch +352}
			fmt.Println("Capital City:", Lëtzebuerg.Capital) // Capital City: Luxembourg-City

			// Prepare the response-data 2
			resp2:= make(map[string]any)
			for k,v:= range r.Form{
				resp[k]=v[0]  
			}

			err=json.Unmarshal(body, &resp2)
			if err!=nil{
				return
			}

			fmt.Println("Unmarshalled json into a map:",resp2)

			// 💡 Access req. details
			fmt.Println("********** ACCESS REQ. DETAILS **********")
			fmt.Println("Body:",r.Body)
			fmt.Println("Form:",r.Form)
			fmt.Println("Header:",r.Header)
			fmt.Println("Context:",r.Context())
			fmt.Println("Content-length:",r.ContentLength)
			fmt.Println("Host:",r.Host)
			fmt.Println("Method:",r.Method)
			fmt.Println("Protocol:",r.Proto)
			fmt.Println("Remote-Addr:",r.RemoteAddr)
			fmt.Println("Req-URI:",r.RequestURI)
			fmt.Println("TLS:",r.TLS)
			fmt.Println("Trailer:",r.Trailer)
			fmt.Println("Transfer-Encoding:",r.TransferEncoding)
			fmt.Println("URL:",r.URL)
			fmt.Println("User-Agent:",r.UserAgent())
			fmt.Println("PORT:",r.URL.Port())
			fmt.Println("Scheme:",r.URL.Scheme)

			
			
// 			********** ACCESS REQ. DETAILS **********
// Body: &{0xc000008018 <nil> <nil> false true {{} {0 0}} true true false 0x886820}
// Form: map[]
// Header: map[Accept:[*  / *] Accept-Encoding:[gzip, deflate, br] Connection:[keep-alive] Content-Length:[128] Content-Type:[application/json] Postman-Token:[5afb66ef-8f20-48a2-8330-f8099b3fbecb] User-Agent:[PostmanRuntime/7.49.1]]
// Context: context.Background.WithValue(net/http context value http-server, *http.Server).WithValue(net/http context value local-addr, [::1]:3000).WithCancel.WithCancel
// Content-length: 128
// Host: localhost:3000
// Method: POST
// Protocol: HTTP/1.1
// Remote-Addr: [::1]:60130
// Req-URI: /teachers
// TLS: <nil>
// Trailer: map[]
// Transfer-Encoding: []
// URL: /teachers
// User-Agent: PostmanRuntime/7.49.1
// PORT:
// Scheme:
			

			
			w.Write([]byte("Hello Post method on Teachers-Route ✅"))
			fmt.Println("Hello Post method on Teachers-Route ✅")
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

	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Students-Route ✅")
		w.Write([]byte("Hello Students-Route ✅"))
		fmt.Println("Students Route ✅", r.Method)
	})

	http.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Execs-Route ✅")
		w.Write([]byte("Hello Execs-Route ✅"))
		fmt.Println("Execs Route ✅", r.Method)
	})

	fmt.Println("Server is running on PORT", PORT,"🟢")
	err:= http.ListenAndServe(PORT,nil)
	if err!=nil{
		log.Fatal("⚠️ERROR. starting the server:",err)
	}
}
```
---

# PATH & QUERY Param(s) ☑️✅

```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// For multiple usecases, use structs (MODELs) instead of simple-maps
type Country struct{
	Name string `json:"country"`
	Capital string `json:"capital"`
	Language string `json:"language"`
	ISD_CODE string `json:"isd-code"`
}

func RootHandler(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Root-Route ✅")
		w.Write([]byte("Hello Root-Route ✅"))
		fmt.Println("Root Route ✅", r.Method)
	}

func TeachersHandler(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Teachers-Route ✅")

		switch r.Method {
				case http.MethodGet:
			//💡 teachers/{id} - Path Params
			//💡 teachers/?key=value&query=value2&sortby=email&sortorder=ASC - Query Params
			fmt.Println("PATH:",r.URL.Path)	//💡 PATH-PARAMS	
			path:=strings.TrimPrefix(r.URL.Path,"/teachers/")
			userID:=strings.TrimSuffix(path,"/")

			fmt.Println("ID is:",userID) // ID is: 123
			fmt.Println("Query Param(s):",r.URL.Query()) // Query Param(s): map[key:[value] query:[value2] sortby:[email] sortorder:[ASC]]
			qryParams:=r.URL.Query()
			sortby:= qryParams.Get("sortby")
			key:= qryParams.Get("key")
			sortorder:= qryParams.Get("sortorder")

			if sortorder == ""{
				sortorder = "DESC"
			}

			fmt.Printf("Sortby: %v, Sort-Order: %v, Key: %v\n", sortby,sortorder,key) // Sortby: email, Sort-Order: ASC, Key: value


			w.Write([]byte("Hello GET method on Teachers-Route ✅"))
			fmt.Println("Hello GET method on Teachers-Route ✅")
			return
				case http.MethodPost:
			//💡 Parse form data (necessary for x-www-form-urlencoded)
			
			err:=r.ParseForm()
			if err!=nil{
				http.Error(w,"ERROR parsing form!",http.StatusBadRequest)
				return
			}

			fmt.Println("FORM:",r.Form)

			// Prepare the response-data
			resp:= make(map[string]any)
			for k,v:= range r.Form{
				resp[k]=v[0]
			}
			fmt.Println("Processed Response Map:",resp)


			//💡 RAW Body { }
			body,err:=io.ReadAll(r.Body)
			if err!=nil{
				return
			}
			r.Body.Close() //⚠️ remember to close the body, bcz, field Body io.ReadCloser!
			fmt.Println("RAW Body:",body)
			

			fmt.Println("Processed Body:", string(body))

			// 💡 UNMARSHAL RAW-Body in JSON format:
			var Lëtzebuerg Country
			err=json.Unmarshal(body, &Lëtzebuerg)
			if err!=nil{
				return
			}
			fmt.Println("Lëtzebuerg:",Lëtzebuerg) // Lëtzebuerg: {Luxembourg Luxembourg-City Lëtzebuergesch +352}
			fmt.Println("Capital City:", Lëtzebuerg.Capital) // Capital City: Luxembourg-City

			// Prepare the response-data 2
			resp2:= make(map[string]any)
			for k,v:= range r.Form{
				resp[k]=v[0]  
			}

			err=json.Unmarshal(body, &resp2)
			if err!=nil{
				return
			}

			fmt.Println("Unmarshalled json into a map:",resp2)

			// 💡 Access req. details
			fmt.Println("********** ACCESS REQ. DETAILS **********")
			fmt.Println("Body:",r.Body)
			fmt.Println("Form:",r.Form)
			fmt.Println("Header:",r.Header)
			fmt.Println("Context:",r.Context())
			fmt.Println("Content-length:",r.ContentLength)
			fmt.Println("Host:",r.Host)
			fmt.Println("Method:",r.Method)
			fmt.Println("Protocol:",r.Proto)
			fmt.Println("Remote-Addr:",r.RemoteAddr)
			fmt.Println("Req-URI:",r.RequestURI)
			fmt.Println("TLS:",r.TLS)
			fmt.Println("Trailer:",r.Trailer)
			fmt.Println("Transfer-Encoding:",r.TransferEncoding)
			fmt.Println("URL:",r.URL)
			fmt.Println("User-Agent:",r.UserAgent())
			fmt.Println("PORT:",r.URL.Port())
			fmt.Println("Scheme:",r.URL.Scheme)

						
			

			
			w.Write([]byte("Hello Post method on Teachers-Route ✅"))
			fmt.Println("Hello Post method on Teachers-Route ✅")
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

func StudentsHandler(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Students-Route ✅")
		w.Write([]byte("Hello Students-Route ✅"))
		fmt.Println("Students Route ✅", r.Method)
	}	

func ExecsHandler(w http.ResponseWriter, r *http.Request){
		//fmt.Fprintf(w, "Hello Execs-Route ✅")
		w.Write([]byte("Hello Execs-Route ✅"))
		fmt.Println("Execs Route ✅", r.Method)
	}

func main() {
	PORT := ":3000"

	http.HandleFunc("/", RootHandler )

	http.HandleFunc("/teachers/", TeachersHandler)

	http.HandleFunc("/students", StudentsHandler)

	http.HandleFunc("/execs", ExecsHandler)

	fmt.Println("Server is running on PORT", PORT,"🟢")
	err:= http.ListenAndServe(PORT,nil)
	if err!=nil{
		log.Fatal("⚠️ERROR. starting the server:",err)
	}
}
```
---
