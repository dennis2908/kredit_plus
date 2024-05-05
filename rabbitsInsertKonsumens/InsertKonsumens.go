package main

import (
	"encoding/json"
	_ "fmt"
	Connection "kredit_plus/connects"
	_ "kredit_plus/routers"
	"log"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"

	"kredit_plus/structs"

	amqp "github.com/rabbitmq/amqp091-go"

	models "kredit_plus/models"
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

func GetData() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
			o := orm.NewOrm()
			o.Using("default")

			ul := &structs.GetKonsumen{}
			json.Unmarshal(d.Body, ul)
			dateBirth, _ := time.Parse("2006-01-02", ul.Date_birth)
			Qry := models.Konsumens{Konsumen: models.Konsumen{Nik: ul.Nik, Full_name: ul.Full_name, Legal_name: ul.Legal_name, Place_birth: ul.Place_birth, Salary: ul.Salary, Foto_ktp: ul.Foto_ktp, Foto_selfie: ul.Foto_selfie}, Date_birth: dateBirth}
			o.Insert(&Qry)

			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
