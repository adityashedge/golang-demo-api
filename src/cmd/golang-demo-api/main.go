package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"
)

const perPage int = 20

var config struct {
	env             string
	session         *mgo.Session
	db              *mgo.Database
	usersCollection *mgo.Collection
}

func init() {
	config.env = "development"
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	config.session = session
	config.db = config.session.DB("golang_demo_api_" + config.env)
	config.usersCollection = config.db.C("users")
}

type ModelErrors map[string][]string

type response struct {
	Message string `json:"message,omitempty"`
	*data   `json:"data,omitempty"`
}

type data struct {
	Total  int `json:"total,omitempty"`
	Users  `json:"users,omitempty"`
	*user  `json:"user,omitempty"`
	Errors ModelErrors `json:"errors,omitempty"`
}

var router = mux.NewRouter().StrictSlash(false)

func main() {
	defer config.session.Close()

	// Negroni Classic has Recovery, Logger and Static.
	// We don't need static file serving in API.
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	log.Println("App initialized...")

	err := http.ListenAndServe(":3000", n)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	return
}
