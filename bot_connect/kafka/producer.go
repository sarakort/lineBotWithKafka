package kafka

import (
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

type Producer struct {
	Con sarama.SyncProducer
}

func newKafkaConfiguration(algorithm string) *sarama.Config {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = 5
	conf.Producer.Return.Successes = true
	// conf.Producer.Return.Errors = true
	conf.ChannelBufferSize = 1
	conf.Version = sarama.V2_2_0_0
	// conf.ClientID = "sasl_scram_client"
	// // conf.Metadata.Full = true
	// conf.Net.SASL.Enable = true
	// // conf.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	// conf.Net.SASL.User = "user"
	// conf.Net.SASL.Password = "JScWkv1Aucah"
	// conf.Net.SASL.Handshake = true
	// if algorithm == "sha512" {
	// 	conf.Net.SASL.SCRAMClient = XDGSCRAMClient{HashGeneratorFcn: SHA512}
	// 	conf.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA512)
	// } else if algorithm == "sha256" {
	// 	conf.Net.SASL.SCRAMClient = DGSCRAMClient{HashGeneratorFcn: SHA256}
	// 	conf.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA256)

	// } else {
	// 	log.Fatalf("invalid SHA algorithm \"%s\": can be either \"sha256\" or \"sha512\"", algorithm)
	// }
	return conf
}

func NewKafkaSyncProducer(brokers string) Producer {
	fmt.Printf("kafka brokers %s\n", brokers)
	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), newKafkaConfiguration("sha256"))
	if err != nil {
		fmt.Printf("Kafka error: %s\n", err)
		os.Exit(-1)
	}
	return Producer{
		Con : producer,
	}
}

func (p Producer) SendByteMsg(topic string, msg []byte) error {
	// fmt.Printf("Topic: %s Message: %+v\n", topic, event)
	// json, err := json.Marshal(event)

	// if err != nil {
	// 	return err
	// }
	fmt.Printf("Topic: %s Json: %s\n", topic, string(msg))

	pmsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(string(msg)),
	}

	partition, offset, err := p.Con.SendMessage(pmsg)

	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in partition %d, offset %d\n", partition, offset)
	return nil
}

func (p Producer)Close(){
	p.Con.Close()
}
