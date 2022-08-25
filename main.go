package main

import (
	"encoding/json"
	"log"
	"wbservice/config"
	"wbservice/controllers/cacheController"
	"wbservice/controllers/httpController"
	"wbservice/controllers/stanController"
	"wbservice/model/order"

	"github.com/nats-io/stan.go"
)

type Application struct {
	stanController  stanController.StanController
	cacheController cacheController.CacheController
	httpController  httpController.HTTPController
	config          config.Config
}

func New() (Application, error) {

	app := new(Application)

	var err error
	app.config, err = config.New("config/config.yml")
	var config *config.Config = &(app.config)
	if err != nil {
		return *app, err
	}

	app.stanController, err = stanController.New(config.Stan.ClusterName, config.Stan.ListenerName, config.Stan.URL)
	if err != nil {
		return *app, err
	}

	err = app.stanController.Subscribe(config.Stan.TopicName, app.handleStanMessage)
	if err != nil {
		return *app, err
	}

	app.cacheController, err = cacheController.New(config.Postgres)
	if err != nil {
		return *app, err
	}

	app.httpController = httpController.New(app.getDataForAPI)
	return *app, nil
}

func (app *Application) Start() error {
	// go publisher.Run(app.config.Stan.ClusterName, "Publisher", app.config.Stan.URL, app.config.Stan.TopicName)
	return app.httpController.Start(app.config.App.Port)
}

func (app *Application) handleStanMessage(msg *stan.Msg) {
	var order order.Order
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Printf("Cannot unmarshal incoming message: %s\n", err.Error())
		return
	}

	uid := order.OrderUID
	if err := app.cacheController.Push(uid, string(msg.Data)); err != nil {
		log.Printf("Cannot push incoming order: %s\n", err.Error())
	}
}

func (app *Application) getDataForAPI(orderUID string) string {
	return app.cacheController.Get(orderUID)
}

func (app *Application) Close() {
	app.cacheController.Close()
	app.stanController.Close()
}

func main() {
	app, err := New()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer app.Close()
	if err := app.Start(); err != nil {
		log.Println(err.Error())
	}
}
