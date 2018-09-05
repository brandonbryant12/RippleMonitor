package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"strconv"
)

func failOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	fmt.Println("Entering PaymentListner")
	conn, err := amqp.Dial(os.Getenv("AMQPConn"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"ripple", // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"client1", // name
		true,      // durable
		false,     // delete when usused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,               // queue name
		"BlockchainListener", // routing key
		"ripple",             // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			//			fmt.Println(string(d.Body))
			var response Response
			json.Unmarshal(d.Body, &response)
			//fmt.Println(response)

			for i := range response.Ledger.Transactions {
				if response.Ledger.Transactions[i].Tx.TransactionType != "Payment" {
					continue
				}
				s, err := strconv.ParseInt(response.Ledger.Transactions[i].Tx.Amount, 10, 64)
				if err != nil {
				}
				amount := s / 1000000
				payment := Payment{
					Currency: "XRP",
					Address:  response.Ledger.Transactions[i].Tx.Destination,
					Amount:   strconv.FormatInt(amount, 10),
					Hash:     response.Ledger.Transactions[i].Hash,
				}

				paymentJSON, err := json.Marshal(payment)
				if err != nil {
					fmt.Println(err)
					return
				}
				err = ch.Publish(
					"ripple",   // exchange
					"payments", // routing key
					false,      // mandatory
					false,      // immediate
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  "text/plain",
						Body:         []byte(paymentJSON),
					})
				fmt.Println(" [x] Sent", payment.String())
				failOnError(err, "Failed to publish a message")

			}
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
