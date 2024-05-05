package main

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	Connection "kredit_plus/connects"
	"kredit_plus/models"
	_ "kredit_plus/routers"
	"log"

	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	Connection.Connects()
}

func main() {

	GetDataUpdate()

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GetDataUpdate() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"updateKonsumens", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
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
			o := orm.NewOrm()
			o.Using("default")

			ul := &models.Konsumens{}
			json.Unmarshal(d.Body, ul)
			// o.Update(&ul)
			_, err := o.Update(ul)
			if err != nil {
				fmt.Println(err.Error())
			}
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
