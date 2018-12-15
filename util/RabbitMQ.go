package util

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type Handler interface {
	Handle(e Event)
}

type Event struct {
	Type    string
	TaskID  int
	Replica int
	Name    string
}

func failOnError(err error, msg string) {
	if err != nil {
		logStr := fmt.Sprintf("%s: %s", msg, err)
		WriteLog(logStr)
	}
}

func createMessage(e Event) string {

	// Create JSON from the instance data.
	// ... Ignore errors.
	b, _ := json.Marshal(e)
	// Convert bytes to string.
	s := string(b)
	return s

}

func Send(qName string, e Event) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		qName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	b, _ := json.Marshal(e)
	// Convert bytes to string.
	body := string(b)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	logstr := fmt.Sprintf(" [x] Sent %s ", body)
	WriteLog(logstr)
	failOnError(err, "Failed to publish a message")
}

func Receive(qName string, handler Handler) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		qName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			logStr := fmt.Sprintf("Received a message: %s", d.Body)
			WriteLog(logStr)
			res := Event{}
			json.Unmarshal([]byte(d.Body), &res)
			handler.Handle(res)
		}
	}()

	WriteLog(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
