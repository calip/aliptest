package models

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	log "github.com/Sirupsen/logrus"
)

type Test struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Title        string        `bson:"title" json:"title"`
	Duration  int        `bson:"duration" json:"duration"`
}

func (test Test) Validate() error {
	if len(test.Title) == 0 {
		log.Errorf("API - Validate - Title is required")
		return errors.New("Validation error")
	}

	if test.Duration == 0 {
		log.Errorf("API - Validate - Duration is required")
		return errors.New("Validation error")
	}

	return nil
}