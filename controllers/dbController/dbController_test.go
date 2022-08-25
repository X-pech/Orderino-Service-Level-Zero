package dbController_test

import (
	"encoding/json"
	"testing"
	"wbservice/config"
	"wbservice/controllers/dbController"
)

func TestDBController(t *testing.T) {
	config, err := config.New("../../config/config.yml")
	if err != nil {
		t.Errorf("Cannot create config object: %s", err.Error())
	}

	pconfig := &config.Postgres
	controller, err := dbController.New(*pconfig)

	if err != nil {
		t.Errorf("Cannot connect to DB: %s", err.Error())
	}

	var testID string = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testObject := struct {
		Test string `json:"test"`
	}{Test: "test"}
	testString, _ := json.Marshal(testObject)
	err = controller.LogToDB(testID, string(testString))

	if err != nil {
		t.Errorf("Cannot Write to DB: %s", err.Error())
	}

	result, err := controller.GetFromDB(testID)
	if err != nil {
		t.Errorf("Cannot get from DB: %s", err.Error())
	}

	if result != string(testString) {
		t.Errorf("Objects are not equal: %s and %s", testString, result)
	}

	_, err = controller.RemoveKey(testID)
	if err != nil {
		t.Errorf("Cannot delete test object from DB: %s", err.Error())
	}

}
