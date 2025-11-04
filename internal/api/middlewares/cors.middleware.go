package middlewares

import (
	"fmt"
	"net/http"
	"slices"
)

// Api is hosted at www.myapi.com
// fontend server is at www.myfrontend.com

// Allowed origins list:
var AllowedOrigins = []string{
	"https://my-origin.url",
	"https://www.myfrontend.com",
	"https://localhost:3000",
}

func CorsMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		origin := r.Header.Get("Origin")

		if IsOriginAllowed(origin){
			w.Header().Set("Access-Control-Allow-Origin",origin)
		 }else{
			http.Error(w, "Not Allowed By CORS ❌", http.StatusForbidden)
			return 
		 }

		fmt.Println(origin)
		w.Header().Set("Access-Control-Allow-Headers","Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods","GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials","true")
		w.Header().Set("Access-Control-Expose-Headers","Authorization")
		w.Header().Set("Access-Control-Max-Age","3600")

		if r.Method == http.MethodOptions{
			// Pre-flight check
			return 
		}

		next.ServeHTTP(w,r)
	})
}

func IsOriginAllowed(origin string)bool{
	return slices.Contains(AllowedOrigins, origin)
}

//! Alternate Way (More verbose.. ⚠️)
// func IsOriginAllowed(origin string)bool{
// 	for _,allowedOrigin:= range AllowedOrigins{
// 		if origin == allowedOrigin{
// 			return true
// 		}
// 	}
// 	return false
// }