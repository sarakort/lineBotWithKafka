package main

import (
	line "line_webhook/channel"
	"line_webhook/kafka"
	"log"
	"os"
	"flag"

	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	consumer := kafka.NewKafkaConsumer(os.Getenv("BOOTSTRAP_SERVERS"),os.Getenv("TOPIC_OUTGOING"))
	// Echo Instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize linebot client
	lineIn := line.NewLineInbound(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"), producer,  os.Getenv("TOPIC_IMCOMING"))
	lineOut:= line.NewLineOutbound(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"), *consumer )

	defer e.Close()
	defer producer.Close()
	defer consumer.Close()

	serviceHandler(e, lineIn)

	go lineOut.ConsumeEvents()
	// Start server listener
	e.Logger.Fatal(e.Start(":8080"))
}

func serviceHandler(e *echo.Echo, lineIn *line.LineInbound) {
	e.GET("/ping", ping)
	e.POST("/webhooks", lineIn.MessageHanlder)
}

//Handler
func ping(c echo.Context) error {
	return c.String(200, "Line boi service is ok!")
}

func main() {
	flag.Parse()
	// loads .env file
	if *useEnv {
		loadEnv()
	}
	startServer()
}
