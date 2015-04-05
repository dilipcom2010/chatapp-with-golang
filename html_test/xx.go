package main

import (
	"github.com/gorilla/context"
	"fmt"
	"time"
	//"runtime"
	"net"
	"net/http"
	//"os"
	//"net/url"
	//"reflect"
	"gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
    "log"
    //"html"
    "github.com/nu7hatch/gouuid"
)

type connection struct {
	ws *websocket.Conn
	send chan []byte
}
type dict map[string]interface{}
var db_coll *mgo.Collection
//var cores int


func main() {
	//cores = runtime.NumCPU()
	//runtime.GOMAXPROCS(cores)
	connect_db()
	http.HandleFunc("/", root)
	http.HandleFunc("/impr", impression)
	http.HandleFunc("/timespent", timespent)
	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/loaderio-88fe9db47315aed418fbb4f1e594b623/", loaderio)
	fmt.Println("listening...")
	err := http.ListenAndServe(":80", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		panic(err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("name", "Dilip Kumar")
	fmt.Fprint(w, "hii")	
}

func redirect(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Cache-Control: no-store, no-cache, must-revalidate, post-check=0, pre-check=0", "max-age=0")

	_, err := r.Cookie("pd_cookie")
	if err != nil {
		u, err := uuid.NewV4()
		//fmt.Println(reflect.TypeOf(u))
		if err == nil {
			cookie := &http.Cookie {
				Name: "pd_cook", 
				Value: u.String(),
				Expires: time.Now().Add(2*365*24*time.Hour),
			}
			http.SetCookie(w, cookie)
			go dump_clicks(r, u)
		}
	}

	 http.Redirect(w,r, r.FormValue("link"), http.StatusMovedPermanently)
}

func impression(w http.ResponseWriter, r *http.Request){
	//runtime.GOMAXPROCS(cores)
	//fmt.Fprint(w, "hii")
	w.Header().Set("Cache-Control:public", "max-age=12000")
	w.Header().Set("Content-Type","text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers","Origin, Accept, Content-Type, X-Requested-With, X-CSRF-Token")
 
	go dump_impression(w, r)
}

func timespent(w http.ResponseWriter, r *http.Request, timespent string) {
	data := make(dict)

	_, err := r.Cookie("pd_cookie")
	if err != nil {
		u, err := uuid.NewV4()
		//fmt.Println(reflect.TypeOf(u))
		if err == nil {
			cookie := &http.Cookie {
				Name: "pd_cook", 
				Value: u.String(),
				Expires: time.Now().Add(2*365*24*time.Hour),
			}
			http.SetCookie(w, cookie)
			data["cookie"] = u.String()
			data["timespent"] = r.FormValue("time")
			data["grouped"] = 0
			data["referer"] = r.Header.Get("Referer")
			//fmt.Println(data)
			err := db_coll.Insert(&data)
			if err != nil {
				log.Fatal(err)
			}
		}
	}		
}


func (c *connection) reader(w http.ResponseWriter, r *http.Request) {
	for {
		_, message, err := c.ws.ReadMessage()
		fmt.Println("success")
		if err != nil {
			break
		}
		go timespent(w, r, string(message))
	}
	c.ws.Close()
	fmt.Println("reader existing.....")
}


func chat(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("regestring connection...")
	var upgrader = &websocket.Upgrader{ReadBufferSize:1024, WriteBufferSize:1024}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("1")
	}
	c := &connection{ws:ws, send:make(chan []byte, 256)}
	c.reader()
	fmt.Fprint(w, "Implemantation is under progress...")
	
}


func loaderio(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "loaderio-88fe9db47315aed418fbb4f1e594b623")
}


func dump_clicks(r *http.Request, u *uuid.UUID) {
	//type dict map[string]interface{}
	data := make(dict)

	var ip string
	fmt.Println("Extracting Datas...")
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		ip = ipProxy
	}
    ip, _, _ = net.SplitHostPort(r.RemoteAddr)

	data["ip"] = ip
	data["link"] = r.FormValue("link")
	data["pub_date"] = r.FormValue("timestamp")
	data["heading"] = r.FormValue("heading")
	data["timestamp"] = time.Now()
	data["cookie"] = u.String()
	data["type"] = "clicks"
	//fmt.Println(data)

	err := db_coll.Insert(&data)
	if err != nil {
		log.Fatal(err)
	}
}


func dump_impression(w http.ResponseWriter, r *http.Request) {
	//type dict map[string]interface{}
	var ip string
	//fmt.Println("Extracting Datas...")
	
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		ip = ipProxy
	}
    ip, _, _ = net.SplitHostPort(r.RemoteAddr)

	data := make(dict)
	data["ip"] = ip
	data["referer"] = r.Header.Get("Referer")
	data["type"] = "impression"
	data["timestamp"] = time.Now()
	//fmt.Println(data["referer"])
	err := db_coll.Insert(&data)
	if err != nil {
		log.Fatal(err)
	}
}


func connect_db() {
	fmt.Println("Connecting Database....")
	session, err := mgo.Dial("119.9.93.32")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")
    //defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    db_coll = session.DB("ctr").C("redirect_stats")
    //fmt.Println(reflect.TypeOf(db.collection))
}
