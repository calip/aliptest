package models

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	log "github.com/Sirupsen/logrus"
)

type Answers struct {
	Text string `json:"text"`
}

type Question struct {
	ID          	bson.ObjectId 	`bson:"_id,omitempty" json:"id,omitempty"`
	Test          	bson.ObjectId 	`bson:"test,omitempty" json:"test,omitempty"`
	Number       	int        	`bson:"number" json:"number"`
	Text       		string        	`bson:"text" json:"text"`
	Answers  		[]Answers       `bson:"answers" json:"answers"`
}

func (question Question) Validate() error {
	if question.Number == 0 {
		log.Errorf("API - Validate - Number is required")
		return errors.New("Validation error")
	}

	if len(question.Text) == 0 {
		log.Errorf("API - Validate - Text is required")
		return errors.New("Validation error")
	}

	return nil
}

