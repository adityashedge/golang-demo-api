package main

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/mgo.v2/bson"
)

type user struct {
	Id                   bson.ObjectId `bson:"_id" json:"id"`
	Name                 string        `bson:"name,omitempty" json:"name,omitempty"`
	Username             string        `bson:"username,omitempty" json:"username,omitempty"`
	Email                string        `bson:"email,omitempty" json:"email,omitempty"`
	Mobile               string        `bson:"mobile,omitempty" json:"mobile,omitempty"`
	Password             string        `bson:"-" json:"-"`
	PasswordConfirmation string        `bson:"-" json:"-"`
	PasswordDigest       string        `bson:"password_digest,omitempty" json:"-"`
	CreatedAt            time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt            time.Time     `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	Errors               ModelErrors   `bson:"-" json:"errors,omitempty"`
}

type newUser struct {
	Name                 *string `bson:"name,omitempty" json:"name,omitempty"`
	Username             *string `bson:"username,omitempty" json:"username,omitempty"`
	Email                *string `bson:"email,omitempty" json:"email,omitempty"`
	Mobile               *string `bson:"mobile,omitempty" json:"mobile,omitempty"`
	Password             *string `bson:"-" json:"password,omitempty"`
	PasswordConfirmation *string `bson:"-" json:"password_confirmation,omitempty"`
}

func (u *user) Create() error {
	u.Id = bson.NewObjectId()

	err := u.generatePasswordDigest()
	if err != nil {
		log.Println(err)
	}

	// before validation callback
	if !u.Valid() {
		return errors.New("User Invalid")
	}
	// after validation callback

	// Update Timestamps
	u.CreatedAt = u.Id.Time()
	u.UpdatedAt = u.CreatedAt

	// before create callback
	err = config.usersCollection.Insert(&u)
	// after create callback
	return err
}

func (u *user) Update(nu newUser) error {
	u.copyFields(nu)

	if nu.Password != nil || nu.PasswordConfirmation != nil {
		err := u.generatePasswordDigest()
		if err != nil {
			log.Println(err)
		}
	}

	// before validation callback
	if !u.Valid() {
		return errors.New("User Invalid")
	}
	// after validation callback

	// Update Timestamps
	u.UpdatedAt = bson.Now()

	// before update callback
	err := config.usersCollection.UpdateId(u.Id, u)
	// after update callback
	return err
}

func (u *user) Valid() bool {
	var userErrors ModelErrors
	if u.Errors != nil {
		userErrors = u.Errors
	} else {
		userErrors = make(ModelErrors)
	}
	// Name, Username and Email are required fields
	// Username and Email must be unique in users collection
	if u.Name == "" {
		userErrors["name"] = append(userErrors["name"], "can't be blank")
	}
	if u.Username == "" {
		userErrors["username"] = append(userErrors["username"], "can't be blank")
	} else {
		n, _ := config.usersCollection.Find(bson.M{"_id": bson.M{"$ne": u.Id}, "username": u.Username}).Count()
		if n > 0 {
			userErrors["username"] = append(userErrors["username"], "is already taken")
		}
	}
	if u.Email == "" {
		userErrors["email"] = append(userErrors["email"], "can't be blank")
	} else {
		n, _ := config.usersCollection.Find(bson.M{"_id": bson.M{"$ne": u.Id}, "email": u.Email}).Count()
		if n > 0 {
			userErrors["email"] = append(userErrors["email"], "is already taken")
		}
	}

	u.Errors = userErrors

	isValid := true
	for _, value := range u.Errors {
		if len(value) > 0 {
			isValid = false
			break
		}
	}
	return isValid
}

func (u *user) copyFields(nu newUser) {
	if nu.Name != nil {
		u.Name = *nu.Name
	}
	if nu.Username != nil {
		u.Username = *nu.Username
	}
	if nu.Email != nil {
		u.Email = *nu.Email
	}
	if nu.Mobile != nil {
		u.Mobile = *nu.Mobile
	}
	if nu.Password != nil {
		u.Password = *nu.Password
	}
	if nu.PasswordConfirmation != nil {
		u.PasswordConfirmation = *nu.PasswordConfirmation
	}
}

func (u *user) generatePasswordDigest() (err error) {
	if u.Errors == nil {
		u.Errors = make(ModelErrors)
	}
	if u.Password == "" {
		u.Errors["password"] = append(u.Errors["password"], "Invalid Password")
		err = errors.New("Invalid Password")
	} else if u.Password != u.PasswordConfirmation {
		u.Errors["password_confirmation"] = append(u.Errors["password_confirmation"], "Password and Password Confirmation do not match")
		err = errors.New("Password and Password Confirmation do not match")
	} else {
		digest, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
		if err == nil {
			u.PasswordDigest = string(digest)
		}
	}
	return
}
