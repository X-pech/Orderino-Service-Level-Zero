package publisher

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nats-io/stan.go"
)

const N = 10000

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GenerateString(c chan string) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < N; i++ {
		uid := make([]rune, 20)
		for j := range uid {
			uid[j] = letters[rand.Intn(len(letters))]
		}
		c <- string(uid)
	}

}

func Run(clusterName string, publisherName string, stanUrl string, topicName string) {
	time.Sleep(5 * time.Second)
	msg := "{\"order_uid\":\"%s\",\"track_number\":\"WBILMTESTTRACK\",\"entry\":\"WBIL\",\"delivery\":{\"name\":\"Test Testov\",\"phone\":\"+9720000000\",\"zip\":\"2639809\",\"city\":\"Kiryat Mozkin\",\"address\":\"Ploshad Mira 15\",\"region\":\"Kraiot\",\"email\":\"test@gmail.com\"},\"payment\":{\"transaction\":\"b563feb7b2b84b6test\",\"request_id\":\"\",\"currency\":\"USD\",\"provider\":\"wbpay\",\"amount\":1817,\"payment_dt\":1637907727,\"bank\":\"alpha\",\"delivery_cost\":1500,\"goods_total\":317,\"custom_fee\":0},\"items\":[{\"chrt_id\":9934930,\"track_number\":\"WBILMTESTTRACK\",\"price\":453,\"rid\":\"ab4219087a764ae0btest\",\"name\":\"Mascaras\",\"sale\":30,\"size\":\"0\",\"total_price\":317,\"nm_id\":2389212,\"brand\":\"Vivienne Sabo\",\"status\":202}],\"locale\":\"en\",\"internal_signature\":\"\",\"customer_id\":\"test\",\"delivery_service\":\"meest\",\"shardkey\":\"9\",\"sm_id\":99,\"date_created\":\"2021-11-26T06:22:19Z\",\"oof_shard\":\"1\"}"
	target := "GET http://localhost:8080/json?order_uid=%s\n\n"
	sc, err := stan.Connect(clusterName, publisherName, stan.NatsURL(stanUrl))
	if err != nil {
		log.Fatal("Not connected")
	}
	c := make(chan string, N)
	go GenerateString(c)
	f, _ := os.OpenFile("./testing/vegeta-test/names.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	for i := 0; i < N; i++ {
		var name string = <-c
		f.WriteString(fmt.Sprintf(target, name))
		err = sc.Publish(topicName, []byte(fmt.Sprintf(msg, name)))
		if err != nil {
			log.Println(err.Error())
		}
	}
}
