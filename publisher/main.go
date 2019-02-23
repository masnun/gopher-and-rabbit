package main

import (
	"encoding/json"
	"github.com/masnun/gopher-and-rabbit"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%: %s", msg, err)
	}

}

func main() {
	conn, err := amqp.Dial(gopher_and_rabbit.Config.AMQPConnectionURL)
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	rand.Seed(time.Now().UnixNano())

	addTask := gopher_and_rabbit.AddTask{Number1: rand.Intn(999), Number2: rand.Intn(999)}
	body, err := json.Marshal(addTask)
	if err != nil {
		handleError(err, "Error encoding JSON")
	}

	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	log.Printf("AddTask: %d+%d", addTask.Number1, addTask.Number2)

}
