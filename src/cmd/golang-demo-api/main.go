package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"

	"gopkg.in/mgo.v2"
)

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
	config.db = config.session.DB("golang_demo_api" + config.env)
	config.usersCollection = config.db.C("users")
}
func main() {
	defer config.session.Close()

	// Negroni Classic has Recovery, Logger and Static.
	// We don't need static file serving in API.
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())

	err := http.ListenAndServe(":3000", n)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}

	log.Println("App initialized...")
	return
}
