package api

import (
	. "../models"
	"crypto/sha256"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"github.com/dgrijalva/jwt-go"
	log "github.com/Sirupsen/logrus"
)

func Login(auth User, r Repository) (Profile, error) {
	var profile Profile
	err := auth.Validate()
	if err != nil {
		return profile, err
	}

	user := User{}
	encrypted := createHash(auth.Password)

	collection := r.DB(DB_NAME).C("users")
	err = collection.Find(bson.M{"email": auth.Email, "password": encrypted}).One(&user)
	if err != nil {
		log.Errorf("API - Login - Fail to find: %s", err)
		return profile, err
	}
	log.Infof("API - Login - Authenticated")

	var token string
	token, err = createToken(auth.Email)
	if err != nil {
		log.Errorf("API - Login - Fail to create token: %s", err)
		return profile, err
	}

	profile.ID = user.ID
	profile.Email = user.Email
	profile.Type = user.Type
	profile.Token = token
	return profile, nil
}

func createHash(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	hmacSecret := []byte("agVrqfQK9CiQZceW")
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func UserFindAll(r Repository) ([]User, error) {
	collection := r.DB(DB_NAME).C("users")
	var users []User
	err := collection.Find(bson.M{"type": "student"}).All(&users)
	return users, err
}

func UserFindById(id string, r Repository) (User, error) {
	var user User
	err := r.DB(DB_NAME).C("users").FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}