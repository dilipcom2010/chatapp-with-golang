package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import fmt



func main() {
	db, err := sql.Open("mysql", "root:@/myapps")
	if err != nil {
		panic(err)
	}
	fmt.Println(db)
}