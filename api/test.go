package api

import (
	log "github.com/Sirupsen/logrus"
	. "../models"
	"gopkg.in/mgo.v2/bson"
)

func TestFindAll(r Repository) ([]Test, error) {
	collection := r.DB(DB_NAME).C("tests")
	var tests []Test
	err := collection.Find(bson.M{}).All(&tests)
	return tests, err
}

func TestFindById(id string, r Repository) (Test, error) {
	var test Test
	err := r.DB(DB_NAME).C("tests").FindId(bson.ObjectIdHex(id)).One(&test)
	return test, err
}

func InsertTest(test Test, r Repository) error {
	err := test.Validate()
	if err != nil {
		return err
	}

	collection := r.DB(DB_NAME).C("tests")
	err = collection.Insert(&test)
	if err != nil {
		log.Errorf("API - Fail to insert: %s", err)
		return err
	}

	log.Infof("API - Test created: %s", test)
	return err
}

func UpdateTest(id string, test Test, r Repository) error {
	err := test.Validate()
	if err != nil {
		return err
	}

	collection := r.DB(DB_NAME).C("tests")
	err = collection.UpdateId(bson.ObjectIdHex(id), &test)
	if err != nil {
		log.Errorf("API - Fail to update: %s", err)
		return err
	}

	log.Infof("API  - Test updated: %s", test)
	return err
}

func DeleteTest(id string, r Repository) error {
	collection := r.DB(DB_NAME).C("tests")
	err := collection.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Errorf("API - Fail to delete ID: %s Error: %s", id, err)
		return err
	}

	log.Infof("API - Test deleted: %s", id)
	return err
}