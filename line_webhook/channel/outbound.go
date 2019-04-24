package channel

import (
	"line_webhook/kafka"
	"encoding/json"
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

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
	var event linebot.Event
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
			
			if err = json.Unmarshal(msgVal, &event ); err != nil{
				fmt.Printf("linebot event unmarshall error: %s\n" ,err)
			}else {
				// logMap = log.(map[string]interface{})
				// replyToken := logMap["replyToken"]
				// fmt.Printf("replyToken %s\n", replyToken)
				if event.Type == linebot.EventTypeMessage {
					switch message := event.Message.(type) {
					case *linebot.TextMessage:
						fmt.Printf("Recived text: %s\n", message.Text)
						messageFromPing := "hello, I am bot." + message.Text
						if err = l.ReplyMessage(event.ReplyToken, messageFromPing); err != nil {
							fmt.Printf("reply message error %s \n", err)
						}
					}
				}
			}	
		}
	}
}
