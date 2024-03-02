package queue

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Message struct {
	Action string `json:"action"`
	Data   []byte `json:"data"`
}

func NewRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://globalhitss:globalhitss@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func PublishMessage(ch *amqp.Channel, queueName string, message *Message) error {
	msgBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBody,
		})
	if err != nil {
		return err
	}

	log.Printf("Sent message to queue: %s", msgBody)
	return nil
}
