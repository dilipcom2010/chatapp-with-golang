package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "fmt"
import "reflect"

type dict map[string]interface{}

func main() {
	//type list []interface{}
	db, err := sql.Open("mysql", "root:dkumar@tcp(localhost:3306)/myapps")
    fmt.Println(reflect.TypeOf(db))
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select id, username, is_active, date_joined from auth_user")
	if err != nil {
		panic(err)
	}
	x := decode(rows)
	fmt.Println(x[0]["username"])
}


func decode(rows *sql.Rows) ([]dict){

	columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    values := make([]sql.RawBytes, len(columns))
    var data []dict


    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    row_no := 0
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }

        var value string
        row := make(dict)
        for i, col := range values {
            if col == nil {
                value = "NULL"
            } else {
                value = string(col)
            }
            row[columns[i]] = value
            //fmt.Println(columns[i], ": ", value)
        }
        data = append(data, row)
    	row_no = row_no + 1
    }

    //fmt.Println(data)
    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    return data
}