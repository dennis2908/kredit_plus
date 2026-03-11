package main

import (
	"context"
	"encoding/json"
	_ "fmt"
	Connection "kredit_plus/connects"

	celery "kredit_plus/celery"
	_ "kredit_plus/routers"
	"log"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"

	"kredit_plus/structs"

	amqp "github.com/rabbitmq/amqp091-go"

	models "kredit_plus/models"

	pusher "kredit_plus/pusherconn"

	"os"
)

func init() {

	Connection.Connects()
}

func main() {

	GetData()

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func failOnError(err error, msg string) string {
	if err != nil {
		return "Error"
	} else {
		return ""
	}
}

func CreateKonsumensMongoMessage(idKonsumen string) string {
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
	msg := failOnError(err, "Failed to connect to RabbitMQ")
	if msg == "Error" {
		return "Error"
	}
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"createKonsumensMongo", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	msg = failOnError(err, "Failed to declare a queue")
	if msg == "Error" {
		return "Error"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(idKonsumen),
		})
	msg = failOnError(err, "Failed to publish a message")
	if msg == "Error" {
		return "Error"
	}
	log.Printf(" [x] Sent %s\n", idKonsumen)

	return "success"

}

func saveData(data []byte) string {
	o := orm.NewOrm()
	o.Using("default")

	ul := &structs.GetKonsumen{}
	json.Unmarshal(data, ul)
	dateBirth, _ := time.Parse("2006-01-02", ul.Date_birth)
	Qry := models.Konsumens{Konsumen: models.Konsumen{Nik: ul.Nik, Full_name: ul.Full_name, Email: ul.Email, Legal_name: ul.Legal_name, Place_birth: ul.Place_birth, Salary: ul.Salary, Foto_ktp: ul.Foto_ktp, Foto_selfie: ul.Foto_selfie}, Date_birth: dateBirth}
	id, _ := o.Insert(&Qry)
	idStr := strconv.Itoa(int(id))
	CreateKonsumensMongoMessage(idStr)
	return "success"
}

func GetData() {
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"createKonsumens", // name
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

			saveData := saveData(d.Body)

			cli, _ := celery.Connect()

			cli.Register("save.data.customers", saveData)

			// start workers (non-blocking call)
			cli.StartWorker()

			// wait for client request
			time.Sleep(4 * time.Second)

			client, _ := pusher.Connect()

			client.Trigger("trigger.load", "konsumens", "all.konsumens")

			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
