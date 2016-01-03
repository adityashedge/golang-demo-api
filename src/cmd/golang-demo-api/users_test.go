package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/mgo.v2/bson"
)

// Index
func TestUsersHandler(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	u := setupUser(t)
	us := Users{u}
	r := response{data: &data{Users: us, Total: 1}}
	expectedResp, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/api/users", nil)
	req.Header.Add("Accept", "application/vnd.demo_app.v1+json")
	req.Header.Add("Content-Type", "application/json")
	actualResp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer actualResp.Body.Close()
	resp, err := ioutil.ReadAll(actualResp.Body)
	if err != nil {
		t.Error(err)
	}
	body := strings.TrimSpace(string(resp))
	if body != string(expectedResp) {
		t.Errorf("expected %s, but got %s", string(expectedResp), body)
	}
	dropAllCollections(t)
}

// Create
func TestCreateUsersHandler(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	u := setupUser(t)
	var b bytes.Buffer
	b.Write([]byte(`{"user":{"name":"Test","username":"test","email":"test@test.com","password":"test123","password_confirmation":"test123"}}`))

	client := &http.Client{}
	req, err := http.NewRequest("POST", ts.URL+"/api/users", &b)
	req.Header.Add("Accept", "application/vnd.demo_app.v1+json")
	req.Header.Add("Content-Type", "application/json")
	actResp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer actResp.Body.Close()
	var resp response
	err = json.NewDecoder(actResp.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	msg := "User successfully created."
	if resp.Message != msg {
		t.Errorf("expected %s, but got %s", msg, resp.Message)
	}

	var nu user
	config.usersCollection.Find(bson.M{"_id": bson.M{"$ne": u.ID}}).One(&nu)
	if nu.Name != "Test" {
		t.Errorf("expected name to be %s, but got %s", "Test", nu.Name)
	}
	if nu.Username != "test" {
		t.Errorf("expected username to be %s, but got %s", "test", nu.Username)
	}
	if nu.Email != "test@test.com" {
		t.Errorf("expected email to be %s, but got %s", "test@test.com", nu.Email)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(nu.PasswordDigest), []byte("test123")); err != nil {
		t.Errorf("expected password to be %s, but got mismatch", "test123")
	}
	totalUsers, _ := config.usersCollection.Find(bson.M{}).Count()
	if totalUsers != 2 {
		t.Errorf("expected %d, but got %d", 2, totalUsers)
	}

	dropAllCollections(t)
}

// Show
func TestShowUserHandler(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	u := setupUser(t)
	r := response{data: &data{user: &u}}
	expResp, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/api/users/"+u.ID.Hex(), nil)
	req.Header.Add("Accept", "application/vnd.demo_app.v1+json")
	req.Header.Add("Content-Type", "application/json")
	actResp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer actResp.Body.Close()
	resp, err := ioutil.ReadAll(actResp.Body)
	if err != nil {
		t.Error(err)
	}
	body := strings.TrimSpace(string(resp))
	if body != string(expResp) {
		t.Errorf("expected %s, but got %s", string(expResp), body)
	}
	dropAllCollections(t)
}

// Edit
func TestEditUserHandler(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	u := setupUser(t)
	r := response{data: &data{user: &u}}
	expResp, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/api/users/"+u.ID.Hex()+"/edit", nil)
	req.Header.Add("Accept", "application/vnd.demo_app.v1+json")
	req.Header.Add("Content-Type", "application/json")
	actResp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer actResp.Body.Close()
	resp, err := ioutil.ReadAll(actResp.Body)
	if err != nil {
		t.Error(err)
	}
	body := strings.TrimSpace(string(resp))
	if body != string(expResp) {
		t.Errorf("expected %s, but got %s", string(expResp), body)
	}
	dropAllCollections(t)
}

// Update
func TestUpdateUsersHandler(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	u := setupUser(t)
	var b bytes.Buffer
	b.Write([]byte(`{"user":{"name":"Test","username":"test","email":"test@test.com","password":"testing123","password_confirmation":"testing123"}}`))

	client := &http.Client{}
	req, err := http.NewRequest("PUT", ts.URL+"/api/users/"+u.ID.Hex(), &b)
	req.Header.Add("Accept", "application/vnd.demo_app.v1+json")
	req.Header.Add("Content-Type", "application/json")
	actResp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer actResp.Body.Close()
	var resp response
	err = json.NewDecoder(actResp.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	msg := "User updated successfully."
	if resp.Message != msg {
		t.Errorf("expected %s, but got %s", msg, resp.Message)
	}

	var nu user
	config.usersCollection.Find(bson.M{"_id": u.ID}).One(&nu)
	if nu.Name != "Test" {
		t.Errorf("expected name to be %s, but got %s", u.Name, nu.Name)
	}
	if nu.Username != "test" {
		t.Errorf("expected username to be %s, but got %s", "test", nu.Username)
	}
	if nu.Email != "test@test.com" {
		t.Errorf("expected email to be %s, but got %s", "test@test.com", nu.Email)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(nu.PasswordDigest), []byte("testing123")); err != nil {
		t.Errorf("expected password to be %s, but got mismatch", "testing123")
	}
	totalUsers, _ := config.usersCollection.Find(bson.M{}).Count()
	if totalUsers != 1 {
		t.Errorf("expected %d, but got %d", 1, totalUsers)
	}

	dropAllCollections(t)
}

// Delete
func TestDeleteUserHandler(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	u := setupUser(t)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", ts.URL+"/api/users/"+u.ID.Hex(), nil)
	req.Header.Add("Accept", "application/vnd.demo_app.v1+json")
	req.Header.Add("Content-Type", "application/json")
	actResp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer actResp.Body.Close()
	var resp response
	err = json.NewDecoder(actResp.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	msg := "User deleted successfully."
	if resp.Message != msg {
		t.Errorf("expected %s, but got %s", msg, resp.Message)
	}

	totalUsers, _ := config.usersCollection.Find(bson.M{}).Count()
	if totalUsers != 0 {
		t.Errorf("expected %d, but got %d", 0, totalUsers)
	}

	dropAllCollections(t)
}
