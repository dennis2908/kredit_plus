package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	_ "fmt"
	models "kredit_plus/models"
	"kredit_plus/structs"
	"log"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	_ "github.com/leekchan/accounting"
	amqp "github.com/rabbitmq/amqp091-go"
	_ "github.com/shopspring/decimal"
)

type KonsumensController struct {
	beego.Controller
}

func (api *KonsumensController) GetAllKonsumens() {
	o := orm.NewOrm()
	o.Using("default")
	var konsumens []models.Konsumens
	num, err := o.QueryTable("konsumens").All(&konsumens)
	if err == nil && num > 0 {
		api.Data["json"] = konsumens
	}
	api.ServeJSON()
}

func (api *KonsumensController) GetKonsumensById() {
	o := orm.NewOrm()
	o.Using("default")
	var konsumens []models.Konsumens
	idInt, _ := strconv.Atoi(api.Ctx.Input.Param(":id"))
	num, err := o.QueryTable("konsumens").Filter("id", idInt).All(&konsumens)
	if err == nil && num > 0 {
		api.Data["json"] = konsumens
	}
	api.ServeJSON()
}

func AllKonsumensCheck(api *KonsumensController) string {
	valid := validation.Validation{}

	frm := api.Ctx.Input.RequestBody
	ul := &structs.GetKonsumen{}
	json.Unmarshal(frm, ul)
	valid.Required(ul.Nik, "Nik")
	valid.Required(ul.Full_name, "Full_name")
	valid.Required(ul.Legal_name, "Legal_name")
	valid.Required(ul.Place_birth, "Place_birth")
	valid.Required(ul.Date_birth, "Date_birth")
	valid.Required(ul.Salary, "Salary")
	valid.Required(ul.Foto_ktp, "Foto_ktp")
	valid.Required(ul.Foto_selfie, "Foto_selfie")

	if valid.HasErrors() {
		// If there are error messages it means the validation didn't pass
		// Print error message
		for _, err := range valid.Errors {
			return err.Key + err.Message
		}
	}

	return ""
}
func failOnError(err error, msg string) string {
	if err != nil {
		return "Error"
	} else {
		return ""
	}
}

func UpdateKonsumensMessage(api *KonsumensController) string {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	msg := failOnError(err, "Failed to connect to RabbitMQ")
	if msg == "Error" {
		return "Error"
	}
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"updateKonsumens", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	msg = failOnError(err, "Failed to declare a queue")
	if msg == "Error" {
		return "Error"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ul := &structs.GetKonsumen{}
	json.Unmarshal(api.Ctx.Input.RequestBody, ul)
	dateBirth, _ := time.Parse("2006-01-02", ul.Date_birth)
	idInt, _ := strconv.Atoi(api.Ctx.Input.Param(":id"))
	// qry := structs.GetKonsumenWID{Id: idInt,Konsumen:models.Konsumen{Nik: ul.Nik, Full_name: ul.Full_name, Legal_name: ul.Legal_name, Place_birth: ul.Place_birth, Salary: ul.Salary, Foto_ktp: ul.Foto_ktp, Foto_selfie: ul.Foto_selfie}, Date_birth: dateBirth}}

	dataSend, err := json.Marshal(models.Konsumens{Id: idInt, Konsumen: models.Konsumen{Nik: ul.Nik, Full_name: ul.Full_name, Legal_name: ul.Legal_name, Place_birth: ul.Place_birth, Salary: ul.Salary, Foto_ktp: ul.Foto_ktp, Foto_selfie: ul.Foto_selfie}, Date_birth: dateBirth})
	if err != nil {
		fmt.Println("error:", err)
	}
	// json.Marshal(Qry)
	body := "send data"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(dataSend),
		})
	msg = failOnError(err, "Failed to publish a message")
	if msg == "Error" {
		return "Error"
	}
	log.Printf(" [x] Sent %s\n", body)

	return "success"

}

func CreateKonsumensMessage(api *KonsumensController) string {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	msg := failOnError(err, "Failed to connect to RabbitMQ")
	if msg == "Error" {
		return "Error"
	}
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"createKonsumens", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	msg = failOnError(err, "Failed to declare a queue")
	if msg == "Error" {
		return "Error"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "send data"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(api.Ctx.Input.RequestBody),
		})
	msg = failOnError(err, "Failed to publish a message")
	if msg == "Error" {
		return "Error"
	}
	log.Printf(" [x] Sent %s\n", body)

	return "success"

}

func (api *KonsumensController) CreateKonsumens() {
	if AllKonsumensCheck(api) != "" {
		api.Data["json"] = AllKonsumensCheck(api)
		api.ServeJSON()
		return
	}

	retData := CreateKonsumensMessage(api)

	if retData == "success" {
		api.Data["json"] = "Successfully insert data"
		api.Ctx.ResponseWriter.WriteHeader(200)
		api.ServeJSON()
	} else {
		api.Ctx.ResponseWriter.WriteHeader(500)
		api.Data["json"] = "Error"

		api.ServeJSON()
	}
}

func (api *KonsumensController) UpdateKonsumens() {
	if AllKonsumensCheck(api) != "" {
		api.Data["json"] = AllKonsumensCheck(api)
		api.ServeJSON()
		return
	}

	retData := UpdateKonsumensMessage(api)

	if retData == "success" {
		api.Data["json"] = "Successfully update data with id " + api.Ctx.Input.Param(":id")
		api.Ctx.ResponseWriter.WriteHeader(200)
		api.ServeJSON()
	} else {
		api.Ctx.ResponseWriter.WriteHeader(500)
		api.Data["json"] = "Error"

		api.ServeJSON()
	}
}
