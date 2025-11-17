package sqlconnect

import (
	"database/sql"
	"fmt"
	_"github.com/go-sql-driver/mysql"
)

func ConnectDB(dbname string) (*sql.DB,error){
	fmt.Println("⏳ Trying to connect to MariaDB...")
	connectionStr:="root:12345@tcp(127.0.0.1:3306)/"+dbname
	db,err:=sql.Open("mysql",connectionStr)
	if err!=nil{
		panic(err)
		// or..
		// return nil,err
	}

	fmt.Println("🛜 Connected to MariaDB!")
	return db,nil
}