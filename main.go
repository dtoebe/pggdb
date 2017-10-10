package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgraph-io/dgraph/client"
	"google.golang.org/grpc"
)

// NewTicket ..
type NewTicket struct {
	Name string
	Desc string
	Tags []string
}

type Ticket struct {
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

var (
	server  *http.Server
	tickets []NewTicket
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Index")
	t, _ := template.ParseFiles("public/base.html", "public/index.html")

	t.ExecuteTemplate(w, "base", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("public/base.html", "public/create.html")

	t.ExecuteTemplate(w, "base", nil)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("public/base.html", "public/admin.html")

	t.ExecuteTemplate(w, "base", &tickets)
}

func save(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, server.Addr+"/", http.StatusSeeOther)
	}

	r.ParseForm()

	log.Println(r.FormValue("customerName"))

	nt := NewTicket{
		Name: r.FormValue("customerName"),
		Desc: r.FormValue("issueDesc"),
		Tags: strings.Split(r.FormValue("tags"), ","),
	}

	tickets = append(tickets, nt)
	http.Redirect(w, r, server.Addr+"/", http.StatusSeeOther)
}

func main() {
	log.Println("Start")
	server = &http.Server{
		Addr: "127.0.0.1:8000",
	}

	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer conn.Close()

	clientDir, err := ioutil.TempDir("", "client_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(clientDir)

	dc := client.NewDgraphClient([]*grpc.ClientConn{conn}, client.DefaultOptions, clientDir)

	req := client.Req{}

	req.SetQuery(`
		{
			tickets(func: eq(type, "ticket")) {
				title
				desc
			}
		}	
	`)

	res, err := dc.Run(context.Background(), &req)
	if err != nil {
		log.Println(err)
	}

	type Root struct {
		Ticket []Ticket `json:"tickets,omitempty"`
	}

	log.Printf("%+v\n", res)

	var t Root
	err = client.Unmarshal(res.N, &t)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%+v\n", t)

	tickets = append(tickets, NewTicket{
		Name: "Dummy customer",
		Desc: "Everything is broken",
		Tags: []string{"Tag1", "Tag2"},
	})

	http.HandleFunc("/save", save)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/", indexHandler)

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
