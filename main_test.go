package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
	"wbservice/model/order"

	"github.com/nats-io/stan.go"
)

const N = 1000

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GenerateString(c chan string) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < N; i++ {
		uid := make([]rune, 30)
		for j := range uid {
			uid[j] = letters[rand.Intn(len(letters))]
		}
		c <- string(uid)
	}

}

func RunService(app *Application) {
	if err := app.Start(); err != nil {
		log.Println(err.Error())
	}
}

func GetExampleOrder(filepath string) (order.Order, error) {
	file, err := os.Open(filepath)
	var result order.Order
	if err != nil {
		return result, err
	}

	s, err := ioutil.ReadAll(file)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(s, &result)

	return result, err
}

func TestOverallIntegration(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Errorf("Cannot create App: %s", err.Error())
	}
	defer app.Close()

	var sample order.Order
	sample, err = GetExampleOrder("./model/order/model.json")
	if err != nil {
		t.Errorf("Cannot find or decode sample order: %s", err.Error())
	}

	sc, err := stan.Connect(app.config.Stan.ClusterName, app.config.Stan.PublisherName, stan.NatsURL(app.config.Stan.URL))
	if err != nil {
		t.Errorf("Cannot connect to nats-streaming: %s", err.Error())
	}

	go RunService(&app)

	c := make(chan string, N)
	go GenerateString(c)

	samples := [N]order.Order{}

	for i := 0; i < N; i++ {
		sample.OrderUID = <-c
		samples[i] = sample
		msg, err := json.Marshal(sample)
		if err != nil {
			t.Errorf("Cannot encode test order to json: %s", err.Error())
		}
		err = sc.Publish(app.config.Stan.TopicName, msg)
		if err != nil {
			t.Errorf("Cannot publish test order: %s", err.Error())
		}
	}

	time.Sleep(3 * time.Second)

	for i := 0; i < N; i++ {
		response, err := http.Get(fmt.Sprintf("http://localhost:%s/json/?order_uid=%s", app.config.App.Port, samples[i].OrderUID))

		if err != nil {
			t.Errorf("Cannot GET order from json-API: %s", err.Error())
		}

		responseBytes, err := ioutil.ReadAll(response.Body)

		if err != nil {
			t.Errorf("Cannot read order from response: %s", err.Error())
		}

		var responseOrder order.Order
		err = json.Unmarshal(responseBytes, &responseOrder)

		if err != nil {
			t.Errorf("Cannot decode recieved order: %s", err.Error())
		}

		if !reflect.DeepEqual(samples[i], responseOrder) {

			b1, _ := json.Marshal(samples[i])
			t.Errorf("Recieved order is different than sent order: \n%s\nvs\n%s\n ", string(b1), string(responseBytes))
		}

	}

}
