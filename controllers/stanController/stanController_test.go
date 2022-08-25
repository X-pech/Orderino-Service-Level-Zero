package stanController_test

import (
	"testing"
	"wbservice/config"
	"wbservice/controllers/stanController"

	"github.com/nats-io/stan.go"
)

type Listener struct {
	msg chan *stan.Msg
}

func newListener() Listener {
	res := new(Listener)
	res.msg = make(chan *stan.Msg, 5)
	return *res
}

func (l *Listener) callback(msg *stan.Msg) {
	l.msg <- msg
}

func TestStanController(t *testing.T) {
	config, err := config.New("../../config/config.yml")
	if err != nil {
		t.Errorf("Cannot create config object: %s", err.Error())
	}

	sconfig := &config.Stan
	controller, err := stanController.New(sconfig.ClusterName, sconfig.ListenerName, sconfig.URL)

	if err != nil {
		t.Errorf("Cannot create nats-streaming controller: %s", err.Error())
	}

	l := newListener()

	topic := "test-topic"
	message := "Hello, world!"

	controller.Subscribe(topic, l.callback)

	publisher, err := stan.Connect(sconfig.ClusterName, sconfig.PublisherName, stan.NatsURL(sconfig.URL))
	if err != nil {
		t.Errorf("Can not connect publisher: %s", err.Error())
	}

	publisher.Publish(topic, []byte(message))

	recievedStanMessage := <-l.msg
	recievedMessage := string(recievedStanMessage.Data)
	if message != recievedMessage {
		t.Errorf("Error: message are not same: %s and %s", message, recievedMessage)
	}

}
