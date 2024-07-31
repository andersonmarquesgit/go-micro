package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"
const MAX_RETRY_CONNECTION = 5

type Config struct {
	RabbitConn *amqp.Connection
}

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		RabbitConn: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func connect() (*amqp.Connection, error) {
	var retryConnection int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready ...")
			retryConnection++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}

		if retryConnection > MAX_RETRY_CONNECTION {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(retryConnection), 2) * float64(time.Second))
		log.Println("backing off ...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil

}
