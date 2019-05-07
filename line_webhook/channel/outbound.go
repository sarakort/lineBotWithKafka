package channel

import (
	"encoding/json"
	"line_webhook/kafka"
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

type BotapiReply struct {
	ReplyToken string `json:"replytoken"`
	Topic string `json:"topic"`
	Confidence string `json:"confidence"`
}

type LineOutbound struct {
	Client *linebot.Client
	Consumer kafka.Consumer
}

func NewLineOutbound(secret string, token string, consumer kafka.Consumer) *LineOutbound{
	return &LineOutbound{
		Client: connect(secret, token),
		Consumer: consumer,
	}
}

func (l *LineOutbound) ReplyMessage(replyToken string, msg string) error {
	fmt.Println(replyToken)
	if _, err := l.Client.ReplyMessage(replyToken, linebot.NewTextMessage(msg)).Do(); err != nil {
		return err
	}
	return nil
}

func (l *LineOutbound)ConsumeEvents(){
	var msgVal []byte
	var err error
	var reply BotapiReply
	// var logMap map[string]interface{}
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
	}()
	  
	for {
		select {
		case err, more := <- l.Consumer.Con.Errors():
			if more {
				fmt.Printf("Kafka consumer error: %s\n", err)
			}
		case msg := <- l.Consumer.Con.Messages():
			l.Consumer.Con.MarkOffset(msg, "")
			msgVal = msg.Value
			fmt.Printf("Topic %s, Processing %s\n",  l.Consumer.Topics, string(msgVal))
			
			if err = json.Unmarshal(msgVal, &reply ); err != nil{
				fmt.Printf("linebot event unmarshall error: %s\n" ,err)
			}else {
				if err = l.ReplyMessage(reply.ReplyToken, reply.Topic); err != nil {
					fmt.Printf("reply message error %s \n", err)
				}
			}	
		}
	}
}
