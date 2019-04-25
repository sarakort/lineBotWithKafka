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
	// fmt.Printf("line secret :  %s, token : %s\n", secret, token)
	return &LineInbound{
		Client: connect(secret, token),
		Producer: producer,
		Topic: topic,
	}
}

func connect(secret string, token string) *linebot.Client {
	bot, err := linebot.New(secret, token)
	if err != nil {
		fmt.Printf("linebot err: %v\n", err)
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
			return c.String(http.StatusBadRequest, linebot.ErrInvalidSignature.Error())
		} else {
			return c.String(http.StatusInternalServerError, "internal server error")
		}
	}

	for _, event := range events {
		fmt.Printf("Event: %+v \n", event)
		if err = l.Producer.SendMsg(l.Topic, event) ; err != nil {
			log.Print(err)
		}
	}
	return c.String(http.StatusOK, "OK!")
}

