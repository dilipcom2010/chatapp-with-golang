package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	//"log"
	"html/template"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

//var homeTempl *template.Template
type dict map[string]interface{}

type connection struct {
	ws *websocket.Conn
	send chan []byte
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/chat", chat)
	fmt.Println("listening...")

	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		panic(err)
	}
}


func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		fmt.Println("success")
		fmt.Println(message)
		c.send <- message
		if err != nil {
			break
		}
		//do something..........
	}
}

func (c *connection) writer() {
	for message := range c.send {
		fmt.Println("hiiii")
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil{
			panic(err)
		}
	}
}





func root(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:dkumar@tcp(localhost:3306)/myapps")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select id, username, is_active, date_joined from auth_user")
	if err != nil {
		panic(err)
	}
	x := decode(rows)
	fmt.Println(x[0]["username"])
	t, _ := template.ParseFiles("home.html")
    t.Execute(w, x)
}


func chat(w http.ResponseWriter, r *http.Request) {
	//var message []byte
	var upgrader = &websocket.Upgrader{ReadBufferSize:1024, WriteBufferSize:1024}
	ws, err := upgrader.Upgrade(w, r, nil)
	//fmt.Println(ws)
	if err != nil {
		//panic(err)
		fmt.Println("1")
		//log.Fatal(err)
	}
	//message[0] = 1
	/*err = ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		//panic(err)
		fmt.Println("2")
		//log.Fatal(err)
	}*/

	c := &connection{ws:ws, send:make(chan []byte, 256)}
	go c.writer()
	c.reader()
	fmt.Fprint(w, "Implemantation is under progress...")
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