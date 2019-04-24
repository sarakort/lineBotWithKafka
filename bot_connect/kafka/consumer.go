package kafka

import (
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

type Consumer struct {
	Con *cluster.Consumer
	Topics string
}

func newKafkaConsumerConfig() *cluster.Config {
	conf := cluster.NewConfig()
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest

	return conf
}

func NewKafkaConsumer(brokers string, topics string) *Consumer {

	consumer, err := cluster.NewConsumer(strings.Split(brokers, ","),
		"bot-consumer",
		strings.Split(topics, ","),
		newKafkaConsumerConfig())

	if err != nil {
		fmt.Printf("Kafka consumer error : %s\n", err)
		os.Exit(-1)
	}

	fmt.Printf("consumer topic %s\n" , topics)

	return &Consumer{
		Con: consumer,
		Topics: topics,
	}
}

func (c Consumer)Close(){
	c.Con.Close()
}