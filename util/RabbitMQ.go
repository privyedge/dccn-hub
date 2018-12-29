package util

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

var RabbitMQHost = "127.0.0.1"

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

func getRabbitMQHost() string {
	host := "amqp://guest:guest@" + RabbitMQHost + ":5672/"
	logStr := fmt.Sprintf("RabbitMQ hostname : %s", host)
	WriteLog(logStr)
	return host
}

//send message to RabbitMQ queue
func Send(qName string, e Event) {

	conn, err := amqp.Dial(getRabbitMQHost())
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

//receive messages from RabbitMQ queue, support failed reconnect function
func Receive(qName string, handler Handler) {
	var rabbitCloseError chan *amqp.Error

	for {
		var conn *amqp.Connection
		var ch *amqp.Channel
		var err error

		for conn == nil || ch == nil {
			conn, err = amqp.Dial(getRabbitMQHost())
			failOnError(err, "Failed to connect to RabbitMQ")

			if conn == nil {
				WriteLog("sleep 30 seconds then retry connecting")
				time.Sleep(30 * time.Second)
				//conn.Close()
			} else {
				ch, err = conn.Channel()
				failOnError(err, "Failed to open a channel")
				if ch == nil {
					WriteLog("sleep 30 seconds then retry connect channel")
					time.Sleep(30 * time.Second)
				}
			}
		}

		//defer ch.Close()

		rabbitCloseError = make(chan *amqp.Error)
		conn.NotifyClose(rabbitCloseError)

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
		msg := <-rabbitCloseError
		//conn.Close()
		logStr := fmt.Sprintf("receive rabbitMQ close messages, error : %s ", msg.Reason)
		WriteLog(logStr)
	}

}
