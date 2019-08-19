package main

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"bot/connect/kafka"
	"bot/connect/channel"
	"log"
	"os"
	"flag"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	useEnv = flag.Bool("env", false, "load .env file flag")
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func startServer() {
	producer := kafka.NewKafkaSyncProducer(os.Getenv("BOOTSTRAP_SERVERS"))
	consumer := kafka.NewKafkaConsumer(os.Getenv("BOOTSTRAP_SERVERS"),os.Getenv("TOPIC_IMCOMING"))
	e := echo.New()

	defer e.Close()
	defer producer.Close()
	defer consumer.Close()

	serviceHandler(e)

	go consumeEvents(consumer, producer , os.Getenv("TOPIC_OUTGOING"), os.Getenv("API_URL"))
	// Start server listener
	e.Logger.Fatal(e.Start(":8080"))
}

func consumeEvents(consumer *kafka.Consumer, producer kafka.Producer, topic string, apiUrl string){
	var msgVal []byte
	var err error
	var event linebot.Event 

	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}


	for {
		select {
		case err, more := <- consumer.Con.Errors():
			if more {
				fmt.Printf("Kafka consumer error: %s\n", err)
			}
		case msg := <- consumer.Con.Messages():
			consumer.Con.MarkOffset(msg, "")
			msgVal = msg.Value
			fmt.Printf("msg in: %s\n", string(msgVal))

			if err = json.Unmarshal(msgVal , &event) ;err !=nil{
				fmt.Println("linebot event unmarshall error")
			}else{
				fmt.Printf("value %v\n", event)
				if event.Type == linebot.EventTypeMessage {
					switch message := event.Message.(type) {
					case *linebot.TextMessage:
						reqUrl  := apiUrl + message.Text
						fmt.Printf("req : %s\n", reqUrl)
						resp, err := client.Get(reqUrl)

						defer resp.Body.Close()

						if err != nil{
							fmt.Printf("call http %s err: %s", apiUrl , err)
							continue
						}

						body, err := ioutil.ReadAll(resp.Body)

						if err != nil {
							fmt.Printf("read body error %s\n", err)
							continue
						}
						var reply channel.BotapiReply

						if err := json.Unmarshal(body, &reply); err != nil {
							fmt.Printf("unmarshall err %s \n", err)
							continue
						}

						reply.ReplyToken = event.ReplyToken

						if err = producer.SendMsg(topic, reply) ; err != nil{
							fmt.Printf("Send to kafka error %s\n", err)
						}

					}
				}

				time.Sleep(1 * time.Second)

			}

		}
	}
}

func serviceHandler(e *echo.Echo) {
	e.GET("/ping", ping)
}

func ping(c echo.Context) error {
	return c.String(200, "boi service is ok!")
}

func main() {
	flag.Parse()
	// loads .env file
	if *useEnv {
		loadEnv()
	}
	startServer()
}
