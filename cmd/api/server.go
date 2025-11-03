package main

import (
	"crypto/tls"
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
	cert:= "cert.pem"
	key:="key.pem"

	mux:= http.NewServeMux()


	mux.HandleFunc("/", RootHandler )

	mux.HandleFunc("/teachers/", TeachersHandler)

	mux.HandleFunc("/students", StudentsHandler)

	mux.HandleFunc("/execs", ExecsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Create custom-server
	server:= &http.Server{
		Addr:PORT,
		Handler: mux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on PORT", PORT,"🟢")
	err:= server.ListenAndServeTLS(cert,key)
	if err!=nil{
		log.Fatal("⚠️ERROR. starting the server:",err)
	}
}