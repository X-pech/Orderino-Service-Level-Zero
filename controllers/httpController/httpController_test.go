package httpController_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"wbservice/config"
	"wbservice/controllers/httpController"
)

type dummyProvider struct {
	dummy string
}

func newDummyProvider(dummy string) dummyProvider {
	return dummyProvider{
		dummy: dummy,
	}
}

func (dp *dummyProvider) callback(_ string) string {
	return dp.dummy
}

func start(controller *httpController.HTTPController, port string) error {
	return controller.Start(port)
}

func TestHttpController(t *testing.T) {
	config, err := config.New("../../config/config.yml")
	if err != nil {
		t.Errorf("Cannot create config object: %s", err.Error())
	}

	dp := newDummyProvider("{\"test\":\"test\"}")

	port := config.App.Port

	controller := httpController.New(dp.callback)

	go start(&controller, port)

	response, err := http.Get(fmt.Sprintf("http://localhost:%s/json?order_uid=a", port))

	if err != nil {
		t.Errorf("Cannot GET from server: %s", err.Error())
	}

	result, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Errorf("Cannot read from response body: %s", err.Error())
	}

	resultString := string(result)

	if dp.dummy != string(resultString) {
		t.Errorf("Objects are not same: %s and %s", dp.dummy, resultString)
	}
}
