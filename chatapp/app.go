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
	"strconv"
	"os"
	"io"
	"github.com/nu7hatch/gouuid"
	"path/filepath"
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

/*  */
type hub struct{
	conns connsDict
	register chan connStrt
	unregister chan connStrt
	broadcast chan []byte
}

var h = hub{
	conns: make(connsDict), /* stores all the sockets of an user */
	register: make(chan connStrt), /* Whenever a new socket is opened it is sent to register via this channel */
	unregister: make(chan connStrt), /* Whenever a socket quits, it is sent to unregister via this channel */
	broadcast:   make(chan []byte), /* All the messages that need broadcast send to this channel to broadcast */
}

//var conns = make(map[string]connection)
var db *sql.DB
var store = sessions.NewCookieStore([]byte("chatapp"))




func (h *hub)run() {
	/* 
		This is a goroutine. Runs forever.
		Responsibility :: 
			Register a new connection
				When an user opens a new tab, a new connection is created.
				Tis connection is send to register WRT that user via register channel.
				This forever loop accepts the connection from the register channel and stores it in WRT the that user id.
			Unregister a connection
				When an user closes a tab then the connection(wrapper of socket and message sending channel) is disconnected.
				So the connection needs to remove.
				This forever loop accepts that disconnected channel and removes it from user id previously it is registered.
			broadcast a message
				It fetches the destination user id from the message.
				gets all the connection of that user and writes the message into their sending channel.

			It also calls a goroutine to write the message into databese.
	*/
	for{
		select{
		case c := <-h.register:
			//h.conns[c.id] = c.connection
			h.conns[c.id] = append(h.conns[c.id], c.connection) /* c.id == id of the user who has opened the new tab */
		
		case c := <-h.unregister:
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
				go dump_msg(msg[0], msg[1], m, true)
			} else{
				go dump_msg(msg[0], msg[1], m, false)
			}
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
	http.HandleFunc("/loadchat", get_prev_talk)
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
	/* connects to the database chatapp with username root and  password fyogi */
	db, _ = sql.Open("mysql", "root:fyogi@tcp(localhost:3306)/chatapp")
}




func (c *connection) reader() {
	/* 
		When the user opens a new tab a new socket has created, then this function starts reading messages from the socket.
		And if it get any message it writes it to broadcast channel to broadcast the message.
	*/
	for {
		_, message, err := c.ws.ReadMessage() /* reading from socket */
		fmt.Println("success")
		if err != nil {
			//it the tab gets closed. then the above ReadMessage() function will not work
			break
		}
		h.broadcast <- message
		//do something..........
	}
	c.ws.Close() /* the above loop will break down once the user closes the tab. Then we need to close the socket */
	fmt.Println("reader existing.....")
}




func (c *connection) writer() {
	/*
		When a user opens a new tab, a new connection is created wrapping two things, the newly created socket and a channel named as send.
		when other connections needs to send message into this tab then they will write message into its send channel.
		here the message will be recieved from this send channel and write it into it's socket.
		the way message is notified on front end.
	*/
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




func root(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}




func update_dp(filename string, id int) {
	/* Updating profile picture of user */
	stmtUpd, err := db.Prepare("update users set selfie=? where id=?")
	if err != nil {
		fmt.Println("unable to prepare updation statement")
	}
	//defer stmtUpd.close()

	_, err = stmtUpd.Exec(filename, id)
	if err != nil {
		fmt.Println("unable to update into db")
	}

}




func validate_dp(buff []byte) bool {
	/* Responsible to validate uploaded file */
	//buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	//_, err = file.Read(buff)

	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}

	filetype := http.DetectContentType(buff)
	fmt.Println(filetype)

	switch filetype {
		case "image/jpeg", "image/jpg":
			return true

		case "image/gif":
			return true

		case "image/png":
			return true

		default:
			return false
		}
}




func room(w http.ResponseWriter, r *http.Request) {
	/*
		When a user inters into chat lobby then this function is called.
		It return an html page that is responsible for the UI of chat lobby.
		If a user uploads profile picture then it is also handled by this function.
		Rteurn chal lobby html page if method is get.
		handle profile picture if method is post.
	*/
	session,_ := store.Get(r, "active")
	if r.Method == "POST" && session.Values["login"]==true{
		file, header, err := r.FormFile("profile-pic")
		if err != nil {
				 fmt.Fprintln(w, err)
				 return
		 }

		/*var file1 = file
		buff := make([]byte, 512) 
		_, err = file1.Read(buff)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		
		if validate_dp(buff) != true {
			fmt.Fprint(w, "please use only jpg/png/gif file format")
			return
		}*/
		//fmt.Println(file1, header.Filename)
		u, _ := uuid.NewV4()
		extn := filepath.Ext(header.Filename)
		if extn != ".jpg" && extn != ".png" && extn != ".gif" && extn != ".jpeg"{
			fmt.Fprint(w, "please use only jpg/png/gif file format")
			return
		}
		filename := "static/images/profile-pic/"+u.String()+extn
		out, err := os.Create(filename)
		if err != nil {
				fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
				return
		}
		//out := "/static/images/profile-pic/"
		_, err = io.Copy(out, file)
		if err != nil {
				fmt.Fprintln(w, err)
		}
		usrid,_ := strconv.Atoi(session.Values["id"].(string))
		filename = "/" + filename
		session.Values["dp"] = filename
		session.Save(r, w)
		update_dp(filename, usrid)
	}



	data := make(dict)

	//session,_ := store.Get(r, "active")
	if session.Values["login"] !=true {
		http.Redirect(w,r, "/login", 302)
	}
	data["id"] = session.Values["id"]
	data["user"] = session.Values["username"]
	data["dp"] = session.Values["dp"]

	
	fmt.Println("hii dilip there is a big bug, root is executin everytime")
	rows, err := db.Query("select Id, FirstName, LastName, selfie, LastActive from users")
	if err != nil {
		panic(err)
	}
	x := decode(rows)
	data["all_users"] = x
	//fmt.Println(data)
	fmt.Println(reflect.TypeOf(x))
	
	t, _ := template.ParseFiles("chat-room.html")
	t.Execute(w, data)
}




func chat(w http.ResponseWriter, r *http.Request) {
	/*
		This is the function called by websocket.
		if a user hits /chat url then this function is called and websocket is created.
		Functionality:
			Starts websocket
			wraps sending channel and socket into a connection
			writes this connection into register channel
			starts a goroutine writer
			starts a goroutine for fetching all unseen messages
			starts a reader...it is forever loop
			when the socket gets close then also reder and writer exist. after that it writes the connection into unregister channel

	*/
	session,err := store.Get(r, "active")
	if err != nil{
		http.Redirect(w,r, "/login", 302)
	}
	
	fmt.Println("regestring connection...")
	var upgrader = &websocket.Upgrader{ReadBufferSize:1024, WriteBufferSize:1024}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("1")
	}
	c := &connection{ws:ws, send:make(chan []byte, 256)}
	h.register <- connStrt{session.Values["id"].(string), c}
	go c.writer()
	go get_unseen_msgs(session.Values["id"].(string))
	c.reader()
	fmt.Fprint(w, "Implemantation is under progress...")
	
	defer func() { h.unregister <- connStrt{session.Values["id"].(string), c} }()
	
}




func get_unseen_msgs(t string) {
	/* 
		Gets all messages from the database whose status is unseen.
		Write these messages into broadcast channel.
	 */
	id,_ := strconv.Atoi(t)
	stmtOut, err := db.Prepare("select id, message from chats where usrTo=? and seen=false")
	if err != nil {
		panic(err)
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(id)
	if err != nil {
		fmt.Println("Error in selecting from db")
	}
	defer rows.Close()

	msgs := decode(rows)
	for _,m := range msgs {
		h.broadcast <- []byte(m["message"].(string))
		go update_msg(m["id"].(string))
	}
	//fmt.Println(msgs)
}




func dump_msg(to string, from string, m []byte, seen bool) {
	/*
		Responsible for writing messages into database
	*/
	usrTo,_ := strconv.Atoi(to)
	//fmt.Println(from)
	usrFrom,_ := strconv.Atoi(from)
	message := string(m)

	//fmt.Println(reflect.TypeOf(usrFrom), reflect.TypeOf(usrTo), reflect.TypeOf(message))
	stmtIns, err := db.Prepare("insert into chats(usrFrom, usrTo, message, seen) values(?, ?, ?, ?)")
	if err != nil {
		fmt.Println("unable to prepare insert statement")
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(usrFrom, usrTo, message, seen)
	if err != nil {
		fmt.Println("unable to insert into db")
	}
}




func update_msg(t string) {
	/*
		It gets the id of message and meake it seen from unseen.
		This happens when the user gets online and seen unseen messages.
	*/
	id,_ := strconv.Atoi(t)
	fmt.Println(id)
	stmtUpd, err := db.Prepare("update chats set seen=true where id=?")
	if err != nil {
		fmt.Println("unable to prepare updation statement")
	}
	//defer stmtUpd.close()

	_, err = stmtUpd.Exec(id)
	if err != nil {
		fmt.Println("unable to update into db")
	}
}




func get_prev_talk(w http.ResponseWriter, r *http.Request) {
	/*
		It fetches all unseen messages.
	*/
	session,_ := store.Get(r, "active")
	if session.Values["login"] !=true {
		http.Redirect(w,r, "/login", 302)
	}


	//id,_ := strconv.Atoi(t)
	offset := r.FormValue("offset")
	limit := r.FormValue("limit")
	u1 := r.FormValue("user1")
	u2 := r.FormValue("user2")
	stmtOut, err := db.Prepare("select message from chats where (usrFrom=? and usrTo=?) or (usrFrom=? and usrTo=?) order by time desc limit ?, ?")
	if err != nil {
		panic(err)
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(u1, u2, u2, u1, offset, limit)
	if err != nil {
		fmt.Println("Error in selecting from db")
	}
	defer rows.Close()

	msgs := decode(rows)

	fmt.Println(msgs)
	//var pp = make([]string, len(msgs))
	l := len(msgs)
	result := ""
	for i := range msgs {
		//xx := msg["message"].(string)
		s := strings.Split(msgs[l-1-i]["message"].(string), ",")
		fmt.Println(s[0], s[1], s[2], s[3:])
		s1 := fmt.Sprint(s[3:])
		s1 = s1[1:len(s1)-1]
		if s[0] == session.Values["id"]{
			result = result + `<div class="msg-box"><div class="bubble-left">`+s1+`&nbsp;&nbsp<span class="time">`+string(s[2])+`</span></div></div>`
		} else {
			result = result + `<div class="msg-box"><div class="bubble-right">`+s1+`&nbsp;&nbsp<span class="time">`+string(s[2])+`</span></div></div>`
		}
		//fmt.Println(reflect.TypeOf(msg["message"]))
	}
	if result == "" {
		result = "No more messages"
	}
	fmt.Fprint(w, result)
}




func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session,_ := store.Get(r, "active")
		if session.Values["login"] == true{
			http.Redirect(w,r, "/room", 302)
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
		session,_ := store.Get(r, "active")
		fmt.Println(session.Values["login"])
		if session.Values["login"] ==true {
			//fmt.Println(session.Values["id"])
			//session.Options = &sessions.Options{MaxAge: -1}
			//session.Save(r, w)
			http.Redirect(w,r, "/room", 302)
		}

		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		email, password := r.FormValue("username"), r.FormValue("password")
		if email != "" {
			rows, err := db.Query("select Id, FirstName, selfie from users where Email='"+email+"' and Password='"+password+"'")
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
				session.Values["dp"] = userdata[0]["selfie"]
				//fmt.Println(session.Values["dp"])
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