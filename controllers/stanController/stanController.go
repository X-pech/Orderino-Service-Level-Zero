package stanController

import (
	"github.com/nats-io/stan.go"
)

type StanController struct {
	clusterID     string
	clientID      string
	natsURL       string
	connection    stan.Conn
	subscriptions []stan.Subscription
}

func New(clusterID string, clientID string, natsURL string) (StanController, error) {
	result := new(StanController)
	result.clientID = clientID
	result.clusterID = clusterID
	result.natsURL = natsURL
	var err error = nil
	result.connection, err = stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	return *result, err
}

func (sc *StanController) Subscribe(topic string, callback stan.MsgHandler) error {
	sub, err := sc.connection.Subscribe(topic, callback)
	if err != nil {
		return err
	}
	sc.subscriptions = append(sc.subscriptions, sub)
	return err
}

func (sc *StanController) Close() error {
	for sub := range sc.subscriptions {
		if err := sc.subscriptions[sub].Unsubscribe(); err != nil {
			return err
		}
	}
	if err := sc.connection.Close(); err != nil {
		return err
	}

	return nil
}
