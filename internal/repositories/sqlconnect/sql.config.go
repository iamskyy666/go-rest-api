package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB,error){
	fmt.Println("Connecting to MariaDB... ⏳")

	// godotenv loads the .env file vars as if they're part of the system OS.
	// mentioning it in the main() is enough
	// err:=godotenv.Load()
	// if err != nil {
	// 	return nil, err
	// }

	user:=os.Getenv("DB_USER")
	password:=os.Getenv("DB_PASSWORD")
	dbname:=os.Getenv("DB_NAME")
	dbport:=os.Getenv("DB_PORT")
	host:=os.Getenv("HOST")
	//connectionStr:="root:12345@tcp(127.0.0.1:3306)/"+dbname
	connectionStr:=fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",user,password,host, dbport,dbname)
	db,err:=sql.Open("mysql",connectionStr)
	if err!=nil{
		panic(err)
		// or..
		// return nil,err
	}

	fmt.Println("Connected to MariaDB! 🛜")
	return db,nil
}