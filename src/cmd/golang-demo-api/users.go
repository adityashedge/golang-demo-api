package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type Users []user

func init() {
	usersRouter := router.Path("/api/users").
		Headers("Content-Type", "application/json", "Accept", "application/vnd.demo_app.v1+json").
		Subrouter()

	usersRouter.Methods("GET").HandlerFunc(usersHandler)
	usersRouter.Methods("POST").HandlerFunc(createUserHandler)

	userRouter := router.PathPrefix("/api/users/{id}").
		Headers("Content-Type", "application/json", "Accept", "application/vnd.demo_app.v1+json").
		Subrouter()

	userRouter.Methods("GET").HandlerFunc(showUserHandler)
	userRouter.Methods("GET").Path("/edit").HandlerFunc(editUserHandler)
	userRouter.Methods("PUT", "PATCH").HandlerFunc(updateUserHandler)
	userRouter.Methods("DELETE").HandlerFunc(deleteUserHandler)
}

// usersHandler returns paginated users in the collection.
// URL: GET /api/users
// HEADERS:
//	"Content-Type": "application/json"
//	"Accept": "application/vnd.botsworth.v1+json"
// PARAMETERS:
//	"page": Current page number(per page 20 records)
func usersHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users Users
	var offset int
	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	if page != 0 {
		offset = (page - 1) * perPage
	}
	config.usersCollection.Find(bson.M{}).Sort("-created_at").Limit(perPage).Skip(offset).All(&users)

	total_users, err := config.usersCollection.Find(bson.M{}).Count()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := response{
		data: &data{
			Total: total_users,
			Users: users,
		},
	}

	userResp, err := json.Marshal(resp)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(userResp)
	}
	return
}

// createUserHandler can be used to add a user.
// URL: POST /api/users
// HEADERS:
//	"Content-Type": "application/json"
//	"Accept": "application/vnd.botsworth.v1+json"
// BODY:
//	user[name]: Name of the user. (required).
//	user[username]: Username of the user. (required)
//	user[email]: Name of the email. (required)
//	user[mobile]: Mobile number of the user.
// EXAMPLE:
//	{
//		"user": {
//			"name": "Aditya Shedge",
//			"username": "aditya",
//			"email": "test@sample.com",
//			"mobile": "9876543210"
//		}
//	}
func createUserHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params struct {
		User struct{ user } `json:"user"`
	}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var resp *response
	u := params.User.user
	var nu = &user{
		Id:       bson.NewObjectId(),
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Mobile:   u.Mobile,
	}

	err = nu.Create()
	if err != nil {
		log.Println(err)
		resp = &response{
			Message: "Unable to save user. Please correct the errors and try again.",
			data: &data{
				Errors: nu.Errors,
				user:   nil,
			},
		}
		w.WriteHeader(422)
	} else {
		resp = &response{Message: "User successfully created.", data: nil}
		w.WriteHeader(http.StatusOK)
	}

	encoder := json.NewEncoder(w)
	if err = encoder.Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

// showUserHandler return the info of a particular user.
// URL: GET /api/users/:id
// HEADERS:
//	"Content-Type": "application/json"
//	"Authorization": "4g27B3m8ZyFRiN8HHvD1u1500yHF9R6G"
// PARAMETERS:
//	"id": ID of the user for which info is to be returned
func showUserHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	var resp *response
	encoder := json.NewEncoder(w)

	if u, err := loadUser(vars["id"]); err != nil {
		log.Println(err)

		resp = &response{Message: "User not found.", data: nil}
		w.WriteHeader(422)
	} else {
		resp = &response{data: &data{user: u}}
		w.WriteHeader(http.StatusOK)
	}
	if err := encoder.Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

// editUserHandler return the info of a particular user for editing.
// URL: GET /api/users/:id/edit
// HEADERS:
//	"Content-Type": "application/json"
//	"Accept": "application/vnd.botsworth.v1+json"
// PARAMETERS:
//	"id": ID of the user for which info is to be returned
func editUserHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	var resp *response
	encoder := json.NewEncoder(w)

	if u, err := loadUser(vars["id"]); err != nil {
		log.Println(err)

		resp = &response{Message: "User not found.", data: nil}
		w.WriteHeader(422)
	} else {
		resp = &response{data: &data{user: u}}
		w.WriteHeader(http.StatusOK)
	}
	if err := encoder.Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

// updateUserHandler can be used to update a particular user.
// URL: PUT|PATCH /api/users/:id
// HEADERS:
//	"Content-Type": "application/json"
//	"Accept": "application/vnd.botsworth.v1+json"
// BODY:
//	user[name]: Name of the user. (required).
//	user[username]: Username of the user. (required)
//	user[email]: Name of the email. (required)
//	user[mobile]: Mobile number of the user.
// EXAMPLE:
//	{
//		"user": {
//			"name": "Aditya Shedge",
//			"username": "aditya",
//			"email": "test@sample.com",
//			"mobile": "9876543210"
//		}
//	}
func updateUserHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params struct {
		// use 'newUser' because if field is left empty intentionally,
		// then json decoder throws error
		User struct{ newUser } `json:"user"`
	}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var resp *response
	nu := params.User.newUser
	u, err := loadUser(mux.Vars(req)["id"])
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(422)
		w.Write([]byte(`{"message": "User not found."}`))
		return
	}

	if err = u.Update(nu); err != nil {
		log.Println(err)
		resp = &response{
			Message: "Unable to update user. Please correct the errors and try again.",
			data:    &data{Errors: u.Errors, user: nil},
		}
		w.WriteHeader(422)
	} else {
		resp = &response{Message: "User updated successfully.", data: nil}
		w.WriteHeader(http.StatusOK)
	}

	encoder := json.NewEncoder(w)
	if err = encoder.Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

// deleteUserHandler Delete a particular user.
// URL: DELETE /api/users/:id
// HEADERS:
//	"Content-Type": "application/json"
//	"Accept": "application/vnd.botsworth.v1+json"
// PARAMETERS:
//	"id": ID of the user to be deleted.
func deleteUserHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	encoder := json.NewEncoder(w)
	var resp *response
	var objectId bson.ObjectId

	// valid bson id
	if valid := bson.IsObjectIdHex(vars["id"]); !valid {
		log.Println(valid)

		resp = &response{Message: "User not found.", data: nil}
		w.WriteHeader(422)
	} else {
		objectId = bson.ObjectIdHex(vars["id"])
		// delete user
		if err := config.usersCollection.RemoveId(objectId); err != nil {
			log.Println(err)

			resp = &response{Message: "Unable to delete user.", data: nil}
			w.WriteHeader(422)
		} else {
			resp = &response{Message: "User deleted successfully.", data: nil}
			w.WriteHeader(http.StatusOK)
		}
	}
	if err := encoder.Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func loadUser(id string) (*user, error) {
	var u user
	var objectId bson.ObjectId

	if valid := bson.IsObjectIdHex(id); !valid {
		return &u, errors.New("Invalid user id.")
	} else {
		objectId = bson.ObjectIdHex(id)
	}

	err := config.usersCollection.FindId(objectId).One(&u)
	return &u, err
}
