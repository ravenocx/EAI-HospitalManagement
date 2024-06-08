package main

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ : %+v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("Failed to open a channel : %+v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"dlx_exchange", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		log.Panicf("Failed to declare an exchange")
	}

	q, err := ch.QueueDeclare(
		"nurse_access", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Panicf("Failed to declare a queue : %+v", err)
	}

	err = ch.QueueBind(
		q.Name,         // queue name
		"",             // routing key
		"dlx_exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed binds an exchange to a queue : %+v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicf("Failed to register a consumer : %+v", err)
	}

	var forever chan struct{}

	go func() {
		for {
			if len(msgs) == 0 {
				log.Printf(" [*] Logging Nurse Access......")
				time.Sleep(5 * time.Minute)
			}
		}
	}()
	go func() {
		for d := range msgs {
			log.Printf("New nurse with access : %s", d.Body)
		}
	}()

	<-forever
}
