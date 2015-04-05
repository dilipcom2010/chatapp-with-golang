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
type connsList []*connection
type connsDict map[string] connsList
type strList [] string

type connection struct {
	ws *websocket.Conn
	send chan []byte
}

type connStrt struct{
	id string
	connection *connection
}


type hub struct{
	conns connsDict
	register chan connStrt
	unregister chan connStrt
	broadcast chan []byte
}
/*type user struct{
	notify chan [] string
	addOnline chan [] string
	offlining chan string
	onlineOf map[string] strList
}
var users = user{
	notify: make(chan []string),
	addOnline: make(chan []string),
	offlining: make(chan string),
	onlineOf: make(map[string] strList),
}*/

var h = hub{
	conns: make(connsDict),
	register: make(chan connStrt),
	unregister: make(chan connStrt),
	broadcast:   make(chan []byte),
}

//var conns = make(map[string]connection)
var db *sql.DB
var store = sessions.NewCookieStore([]byte("chatapp"))


func (h *hub)run() {
	for{
		select{
		case c := <-h.register:
			//h.conns[c.id] = c.connection
			h.conns[c.id] = append(h.conns[c.id], c.connection)
		
		case c := <-h.unregister:
			//if h.conns[c.id] == c.connection {
			//	delete(h.conns, c.id)
			//	close(c.connection.send)
			//}
			//var temp connsList
			fmt.Println("deleting connection......")
			count := 0
			for _,connec := range h.conns[c.id] {
				if connec == c.connection{
					break
				}
				count = count + 1
			}
			x := h.conns[c.id]
			h.conns[c.id] = append(x[0:count], x[count+1:]...)
			//fmt.Println("conn len==", len(h.conns[c.id]))
			/*if len(h.conns[c.id]) <= 0 {
				//h.offline <- c.id
				go go_offline(c.id)
				users.offlining <- c.id
				fmt.Println("offlining...")
			}*/
			close(c.connection.send)
		
		case m := <-h.broadcast:
			fmt.Println("broadcasting....")
			msg := strings.Split(string(m), ",")
			fmt.Println(msg[0])
			if conn, ok := h.conns[msg[0]]; ok {
				fmt.Println(msg)
				//conn.send <- m
				for _,x := range conn{
					select {
						case x.send <- m:
						default:
							fmt.Println("buffer is full")
					}
				}
			}

		/*case off := <-h.offline:
			fmt.Println("offlining....")
			msg := strings.Split(string(m), ",")
			fmt.Println(msg[0])
			if conn, ok := h.conns[msg[0]]; ok {
				fmt.Println(msg)
				//conn.send <- m
				for _,x := range conn{
					select {
						case x.send <- m:
						default:
							fmt.Println("buffer is full")
					}
				}
			}*/
		}
	}	
}


func main() {
	connectdb()
	go h.run()
	//go users.startnotify()
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })
	http.HandleFunc("/", root)
	http.HandleFunc("/room", room)
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	fmt.Println("listening...")

	err := http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
	if err != nil{
		panic(err)
	}
}



func connectdb() {
	db, _ = sql.Open("mysql", "root:dkumar@tcp(localhost:3306)/chatapp")
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
		//fmt.Println(string(message))
		//--c.send <- message
		if err != nil {
			break
		}
		h.broadcast <- message
		//do something..........
	}
	c.ws.Close()
	fmt.Println("reader existing.....")
}

func (c *connection) writer() {
	for message := range c.send {
		fmt.Println("message====", string(message))
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil{
			panic(err)
		}
	}
	c.ws.Close()
	fmt.Println("writer existing.....")
}

/*func (users *user) startnotify() {
	//var xyz []byte
	for{
		select{
		case n := <-users.notify:
			fmt.Println("working...", n[0], n[1])
			//xyz := []byte(n[0])
			//xyz1 := []byte(n[1])
			//xyz2 := append(xyz, xyz1)
			//n[1] = "-1"
			xyz := strings.Join(n, ",")
			//fmt.Println(xyz)
			h.broadcast <- []byte(xyz)

		case a := <-users.addOnline:
			users.onlineOf[a[0]] = append(users.onlineOf[a[0]], a[1])
			//fmt.Println("online list... ",users.onlineOf["1"])

		case id := <-users.offlining:
			fmt.Println("offline processing....")
			for _,usr := range users.onlineOf[id]{
				fmt.Println("notifying...", usr)
				users.notify <- []string{usr, id, "-2"}
			}
		}
	}
}*/


func root(w http.ResponseWriter, r *http.Request) {
}


func room(w http.ResponseWriter, r *http.Request) {
	data := make(dict)

	session,_ := store.Get(r, "active")
	if session.Values["login"] !=true {
		http.Redirect(w,r, "/login", 302)
	}
	data["id"] = session.Values["id"]
	data["user"] = session.Values["username"]

	fmt.Println("hii dilip there is a big bug, root is executin everytime")
	rows, err := db.Query("select Id, FirstName, LastName, selfie, LastActive from users")
	if err != nil {
		panic(err)
	}
	x := decode(rows)
	data["all_users"] = x
	//fmt.Println(data)
	fmt.Println(reflect.TypeOf(x))
	
	t, _ := template.ParseFiles("home1.html")
    t.Execute(w, data)
}

/*func getonlineusers(au string) {
	rows, err := db.Query("select Id from users where online=true")
	if err != nil {
		panic(err)
	}
	x := decode(rows)
	for _,oui := range x {
		//fmt.Println("hiiiii ", aui)
		users.onlineOf[au] = append(users.onlineOf[au], oui["Id"].(string))
		users.notify <- []string{oui["Id"].(string), au, "-1"}//for notifying all of his friend
		users.notify <- []string{au, oui["Id"].(string), "-1"}//to notify himself
		//users.addOnline <- []string{oui["Id"].(string), au}
		//users.addOnline <- []string{au, oui["Id"].(string)}
		//notification_of_being_online_to(au["Id"].(string), au)
	}
	fmt.Println("online list... ",users.onlineOf[au])
}
func go_online(au string) {
	_, err := db.Query("update users set online=true where id="+au)
	if err != nil {
		panic(err)
	}
}
func go_offline(au string) {
	_, err := db.Query("update users set online=false where id="+au)
	if err != nil {
		panic(err)
	}
}*/

func chat(w http.ResponseWriter, r *http.Request) {
	//var message []byte
	session,_ := store.Get(r, "active")
	//id := session.Values["id"].(string)
	
	fmt.Println("regestring connection...")
	var upgrader = &websocket.Upgrader{ReadBufferSize:1024, WriteBufferSize:1024}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("1")
	}
	c := &connection{ws:ws, send:make(chan []byte, 256)}
	h.register <- connStrt{session.Values["id"].(string), c}
	go c.writer()
	/*if session.Values["get_ol9"] == true {
		 go getonlineusers(session.Values["id"].(string))
		 go go_online(session.Values["id"].(string))	
	}*/
	c.reader()
	fmt.Fprint(w, "Implemantation is under progress...")
	
	defer func() { h.unregister <- connStrt{session.Values["id"].(string), c} }()
	
}



func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session,_ := store.Get(r, "active")
		if session.Values["login"] == true{
			http.Redirect(w,r, "/", 302)
		}
		t,_ := template.ParseFiles("signup.html")
		t.Execute(w, nil)
	} else {
		first_name := r.FormValue("FirstName")
		last_name := r.FormValue("LastName")
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Println(first_name, last_name, email, password, reflect.TypeOf(password))
		_, err := db.Query("insert into users(FirstName, LastName, Email, Password) values('"+first_name+"', '"+last_name+"', '"+email+"', '"+password+"')")
		if err != nil {
			fmt.Println("unable to insert into db")
		}
		http.Redirect(w,r, "/login", 302)
	}
	
}


func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//clearSession(w)
		session,_ := store.Get(r, "active")
		fmt.Println(session.Values["login"])
		if session.Values["login"] ==true {
			fmt.Println(session.Values["id"])
			session.Options = &sessions.Options{MaxAge: -1}
			session.Save(r, w)
			http.Redirect(w,r, "/room", 302)
		}

		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		email, password := r.FormValue("username"), r.FormValue("password")
		if email != "" {
			rows, err := db.Query("select Id, FirstName from users where Email='"+email+"' and Password='"+password+"'")
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
				session.Values["username"] = userdata[0]["FirstName"]
				session.Values["id"] = userdata[0]["Id"]
				session.Values["get_ol9"] = true
				session.Save(r, w)
				//http.Redirect(w,r, "/", http.StatusMovedPermanently)
				http.Redirect(w,r, "/room", 302)
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