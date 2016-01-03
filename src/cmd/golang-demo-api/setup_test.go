package main

import "testing"

func init() {
	config.env = "test"
	config.db = config.session.DB("golang_demo_api_" + config.env)
	config.usersCollection = config.db.C("users")
}

func dropAllCollections(t *testing.T) {
	err := config.usersCollection.DropCollection()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func setupUser(t *testing.T) user {
	u := user{
		Name:                 "Test User",
		Username:             "test_user",
		Email:                "test@sample.com",
		Mobile:               "9876543210",
		Password:             "test123#",
		PasswordConfirmation: "test123#",
	}
	err := u.Create()
	if err != nil {
		t.Fatal("Unable to create user: ", err, u.Errors)
	}
	return u
}
