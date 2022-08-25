package cacheController_test

import (
	"testing"
	"wbservice/config"
	"wbservice/controllers/cacheController"
)

func TestCacheController(t *testing.T) {
	config, err := config.New("../../config/config.yml")
	if err != nil {
		t.Errorf("Cannot create config object: %s", err.Error())
	}

	pconfig := &config.Postgres
	controller, err := cacheController.New(*pconfig)

	if err != nil {
		t.Errorf("Cannot create Cache Controller: %s", err.Error())
	}

	testID := "aaaaaa"
	testObj := "{\"test\":\"test\"}"
	err = controller.Push(testID, testObj)
	if err != nil {
		t.Errorf("Cannot Push object: %s", err.Error())
	}

	resultString := controller.Get(testID)
	if resultString != testObj {
		t.Errorf("Objects are not equal: %s and %s", testObj, resultString)
	}

	controller.Erase(testID)

}
