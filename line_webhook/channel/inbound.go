package channel

import (
	"line_webhook/kafka"
	"context"
	"fmt"
	"net/http"
	"os"
	"log"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineInbound struct {
	Client *linebot.Client
	Producer kafka.Producer
	Topic string
}

func NewLineInbound(secret string, token string, producer kafka.Producer, topic string) *LineInbound {
	return &LineInbound{
		Client: connect(secret, token),
		Producer: producer,
		Topic: topic,
	}
}

func connect(secret string, token string) *linebot.Client {
	bot, err := linebot.New(secret, token)
	if err != nil {
		fmt.Printf("linebot err: %v", err)
		os.Exit(-1)
	}
	return bot
}

func (l *LineInbound) MessageHanlder(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	events, err := l.Client.ParseRequest(c.Request())
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(400, linebot.ErrInvalidSignature.Error())
		} else {
			c.String(500, "internal")
		}
	}

	for _, event := range events {
		fmt.Printf("Event: %+v \n", event)
		if err = l.Producer.SendMsg(l.Topic, event) ; err != nil {
			log.Print(err)
		}
	}
	return c.JSON(http.StatusOK, "OK!")
}

