package api

import (
	log "github.com/Sirupsen/logrus"
	. "../models"
	"gopkg.in/mgo.v2/bson"
)

func QuestionFindAll(r Repository) ([]Question, error) {
	collection := r.DB(DB_NAME).C("questions")
	var questions []Question
	err := collection.Find(bson.M{}).All(&questions)
	return questions, err
}

func QuestionFindById(id string, r Repository) (Question, error) {
	var question Question
	err := r.DB(DB_NAME).C("questions").FindId(bson.ObjectIdHex(id)).One(&question)
	return question, err
}

func InsertQuestion(question Question, r Repository) error {
	err := question.Validate()
	if err != nil {
		return err
	}

	collection := r.DB(DB_NAME).C("questions")
	err = collection.Insert(&question)
	if err != nil {
		log.Errorf("API - Fail to insert: %s", err)
		return err
	}

	log.Infof("API - Question created: %s", question)
	return err
}

func UpdateQuestion(id string, question Question, r Repository) error {
	err := question.Validate()
	if err != nil {
		return err
	}

	collection := r.DB(DB_NAME).C("questions")
	err = collection.UpdateId(bson.ObjectIdHex(id), &question)
	if err != nil {
		log.Errorf("API - Fail to update: %s", err)
		return err
	}

	log.Infof("API  - Question updated: %s", question)
	return err
}

func DeleteQuestion(id string, r Repository) error {
	collection := r.DB(DB_NAME).C("questions")
	err := collection.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Errorf("API - Fail to delete ID: %s Error: %s", id, err)
		return err
	}

	log.Infof("API - Question deleted: %s", id)
	return err
}