package main

import (
	"context"
	"fmt"
	_ "fmt"
	mongoconn "kredit_plus/mongoconn"
	_ "kredit_plus/routers"
	"log"

	_ "github.com/lib/pq"

	"kredit_plus/structs"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	GetData()

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GetData() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"updateKonsumensMongo", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {

		for d := range msgs {
			post := structs.InsertKonsumen{
				IdKonsumen: string(d.Body),
				Operation:  "mongo update data konsumen",
			}
			db, err := mongoconn.Connect()
			if err != nil {
				log.Fatal(err.Error())
			}

			var ctx = context.TODO()

			// Insert ke database
			_, errx := db.Collection("konsumens").InsertOne(ctx, post)

			// Handle error
			if errx != nil {
				fmt.Printf("an error ocurred when connect to mongoDB : %v", err)
				panic(err)
			}

			fmt.Println("Proses update berhasil...")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
