package models

import (
	"gopkg.in/mgo.v2/bson"
	log "github.com/Sirupsen/logrus"
	"errors"
)

type User struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Type string `json:"type"`
}

type Profile struct {
	ID        	bson.ObjectId `bson:"_id" json:"id"`
	Email      	string
	Type		string
	Token		string
}

func (user User) Validate() error {
	if len(user.Email) == 0 {
		log.Errorf("API - Validate - Email is required")
		return errors.New("Validation error")
	}

	if len(user.Password) == 0 {
		log.Errorf("API - Validate - Password is required")
		return errors.New("Validation error")
	}
	return nil
}