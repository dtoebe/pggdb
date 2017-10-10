package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc"
)

// NewTicket ..
type NewTicket struct {
	Name string
	Desc string
	Tags []string
}

var (
	server  *http.Server
	tickets []NewTicket
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
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
	server = &http.Server{
		Addr: "127.0.0.1:8000",
	}

	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	tickets = append(tickets, NewTicket{
		Name: "Dummy customer",
		Desc: "Everything is broken",
		Tags: []string{"Tag1", "Tag2"},
	})

	http.HandleFunc("/save", save)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/", indexHandler)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
