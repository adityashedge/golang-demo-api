package main

import (
	"log"

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
	log.Println("App initialized...")
	return
}
