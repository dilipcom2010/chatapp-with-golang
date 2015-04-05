package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
	"net/http"
	//"log"
	"html/template"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"reflect"
	"strings"
)

//var homeTempl *template.Template
type dict map[string]interface{}

type connection struct {
	ws *websocket.Conn
	send chan []byte
}

//type hub struct{
//	conns map[string]connection
//}

//var h = hub{
//	conns: make(map[string]connection),
//}

var dd = 10
//var conns = make(map[string]connection)
var db *sql.DB
var store = sessions.NewCookieStore([]byte("chatapp"))


func main() {
	fmt.Println(dd)
	dd=98
	connectdb()
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })
	http.HandleFunc("/", root)
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	fmt.Println("listening...")

	err := http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
	if err != nil{
		panic(err)
	}
}



func connectdb() {
	db, _ = sql.Open("mysql", "root:dkumar@tcp(localhost:3306)/myapps")
	//fmt.Println(reflect.TypeOf(err))
	//fmt.Println(db)
	//if err != nil {
	//	panic(err)
	//}
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		fmt.Println("success")
		m := strings.Split(string(message), ",")
		fmt.Println(m[2:])
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
	data := make(dict)

	session,_ := store.Get(r, "active")
	if session.Values["login"] !=true {
		http.Redirect(w,r, "/login", 302)
	}
	data["id"] = session.Values["id"]
	data["user"] = session.Values["username"]

	rows, err := db.Query("select id, username, is_active, date_joined from auth_user")
	if err != nil {
		panic(err)
	}
	x := decode(rows)
	data["all_users"] = x
	//fmt.Println(data)
	fmt.Println(reflect.TypeOf(x))
	t, _ := template.ParseFiles("home.html")
    t.Execute(w, data)
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

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//clearSession(w)
		session,_ := store.Get(r, "active")
		fmt.Println(session.Values["login"])
		if session.Values["login"] ==true {
			fmt.Println(session.Values["id"])
			//session.Options = &sessions.Options{MaxAge: -1}
			//session.Save(r, w)
			http.Redirect(w,r, "/", 302)
		}

		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		username, _ := r.FormValue("username"), r.FormValue("password")
		if username != "" {
			rows, err := db.Query("select id, username from auth_user where username='"+username+"'")
			if err != nil {
				//panic(err)
				fmt.Println(rows)
			}
			userdata := decode(rows)
			if userdata == nil {
				t, _ := template.ParseFiles("login.html")
				t.Execute(w, nil)
			} else {
				session, _ := store.Get(r, "active")
				//fmt.Println(session.Values["foo"])
				session.Values["login"] = true
				session.Values["username"] = username
				session.Values["id"] = userdata[0]["id"]
				session.Save(r, w)
				//http.Redirect(w,r, "/", http.StatusMovedPermanently)
				http.Redirect(w,r, "/", 302)
			}
		} else {
			t, _ := template.ParseFiles("login.html")
			t.Execute(w, nil)
		}
	}
}


func logout(w http.ResponseWriter, r *http.Request) {
	session,_ := store.Get(r, "active")
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(r, w)
	http.Redirect(w,r, "/login", 302)
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