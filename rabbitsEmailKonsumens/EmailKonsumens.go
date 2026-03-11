package main

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	Connection "kredit_plus/connects"
	_ "kredit_plus/routers"
	"log"
	"net/smtp"
	"os"

	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"

	"kredit_plus/structs"

	amqp "github.com/rabbitmq/amqp091-go"
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
		"sendEmailToKonsumens", // name
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
			log.Printf("Received a message: %s", d.Body)
			o := orm.NewOrm()
			o.Using("default")

			ul := &structs.GetKonsumen{}
			json.Unmarshal(d.Body, ul)

			// SMTP configuration
			username := os.Getenv("SMTP_USERNAME")
			password := os.Getenv("SMTP_PASSWORD")
			host := os.Getenv("SMTP_HOST")
			port := os.Getenv("SMTP_PORT")

			// Subject and body
			subject := "Welcome, " + ul.Full_name
			body := "Hi, " + ul.Full_name + ". Welcome to our company"

			// Sender and receiver
			from := "michaeldenniseldima@gmail.com"
			to := []string{
				ul.Email,
			}

			// Build the message
			message := fmt.Sprintf("From: %s\r\n", from)
			message += fmt.Sprintf("To: %s\r\n", to)
			message += fmt.Sprintf("Subject: %s\r\n", subject)
			message += fmt.Sprintf("\r\n%s\r\n", body)

			// Authentication.
			auth := smtp.PlainAuth("", username, password, host)

			// Send email
			err := smtp.SendMail(host+":"+port, auth, from, to, []byte(message))
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Email sent successfully.")

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
