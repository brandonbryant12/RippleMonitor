package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

var (
	outfile, _ = os.Create("/data/stellarPayments.log") // update path for your needs
	l          = log.New(outfile, "", 0)
)

func main() {
	l.Println("Entering PaymentListner")
	conn, err := amqp.Dial(os.Getenv("AMQPConn"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"stellar", // name
		"direct",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
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
		"stellar",            // exchange
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
			fmt.Println(string(d.Body))
			var response Response
			json.Unmarshal(d.Body, &response)
			if response.Ledger.Accepted != "true" {
				continue
			}
			fmt.Println(response)

			/*			txHash := response.Ledger.Transactions.Hash

						for i := range response.Ledger.Transactions.Tx {
							if response.Ledger.Transactions.Tx[i].TransactionType != "Payment" {
								continue
							}
							payment := Payment{
								Currency: "XLM",
								Address:  response.Ledger.Transactions.Tx[i].Destination,
								Amount:   response.Ledger.Transactions.Tx[i].Amount,
								Hash:     txHash,
							}

							paymentJSON, err := json.Marshal(payment)
							if err != nil {
								fmt.Println(err)
								return
							}
							err = ch.Publish(
								"stellar",  // exchange
								"payments", // routing key
								false,      // mandatory
								false,      // immediate
								amqp.Publishing{
									DeliveryMode: amqp.Persistent,
									ContentType:  "text/plain",
									Body:         []byte(paymentJSON),
								})
							fmt.Printf(" [x] Sent %s", payment.String())
							failOnError(err, "Failed to publish a message")

						}*/
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
