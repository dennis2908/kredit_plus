package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	_ "fmt"
	models "kredit_plus/models"
	mongoconn "kredit_plus/mongoconn"
	"kredit_plus/structs"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	_ "github.com/leekchan/accounting"
	amqp "github.com/rabbitmq/amqp091-go"
	_ "github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
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

func (api *KonsumensController) GetAllKonsumensMongoInsert() {
	db, err := mongoconn.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	var ctx = context.TODO()

	csr, err := db.Collection("konsumens").Find(ctx, bson.M{"content": "mongo insert data konsumen"})

	result := make([]structs.InsertKonsumen, 0)
	for csr.Next(ctx) {
		var row structs.InsertKonsumen
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}
	if err == nil {
		api.Data["json"] = result
	}

	api.ServeJSON()
}

func (api *KonsumensController) GetAllKonsumensMongoUpdate() {
	db, err := mongoconn.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	var ctx = context.TODO()

	csr, err := db.Collection("konsumens").Find(ctx, bson.M{"content": "mongo update data konsumen"})

	result := make([]structs.InsertKonsumen, 0)
	for csr.Next(ctx) {
		var row structs.InsertKonsumen
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}
	if err == nil {
		api.Data["json"] = result
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
	valid.Required(ul.Email, "Email")
	valid.Email(ul.Email, "Email")
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

	dataSend, err := json.Marshal(models.Konsumens{Id: idInt, Konsumen: models.Konsumen{Nik: ul.Nik, Full_name: ul.Full_name, Email: ul.Email, Legal_name: ul.Legal_name, Place_birth: ul.Place_birth, Salary: ul.Salary, Foto_ktp: ul.Foto_ktp, Foto_selfie: ul.Foto_selfie}, Date_birth: dateBirth})
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

func EmailKonsumensMessage(api *KonsumensController) string {
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
		"sendEmailToKonsumens", // name
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

func (api *KonsumensController) ExcelKonsumens() {

	o := orm.NewOrm()
	o.Using("default")
	var konsumens []models.Konsumens
	num, err := o.QueryTable("konsumens").All(&konsumens)
	if err == nil && num > 0 {
		api.Data["json"] = konsumens
	}
	xlsx := excelize.NewFile()
	sheetName := "Sheet1"

	xlsx.SetSheetName("Sheet1", sheetName)

	// Add headers
	xlsx.SetCellValue(sheetName, "A1", "Email")
	xlsx.SetCellValue(sheetName, "B1", "Nik")
	xlsx.SetCellValue(sheetName, "C1", "Full_name")
	xlsx.SetCellValue(sheetName, "D1", "Legal_name")
	xlsx.SetCellValue(sheetName, "E1", "Place_birth")
	xlsx.SetCellValue(sheetName, "F1", "Date_birth")
	xlsx.SetCellValue(sheetName, "G1", "Salary")
	xlsx.SetCellValue(sheetName, "H1", "Foto_ktp")
	xlsx.SetCellValue(sheetName, "I1", "Foto_selfie")
	// Create a new sheet.
	rowIndex := 2
	for _, data := range konsumens {
		xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), data.Email)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), data.Nik)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), data.Full_name)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), data.Legal_name)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), data.Place_birth)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIndex), data.Date_birth)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIndex), data.Salary)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("H%d", rowIndex), data.Foto_ktp)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("I%d", rowIndex), data.Foto_selfie)

		rowIndex++
	}

	filename := "static/excel/" + api.Ctx.Input.Param(":filename") + ".xlsx"

	if err := xlsx.SaveAs(filename); err != nil {
		log.Fatal(err)
	}
}

func (api *KonsumensController) CreateKonsumens() {
	if AllKonsumensCheck(api) != "" {
		api.Data["json"] = AllKonsumensCheck(api)
		api.ServeJSON()
		return
	}

	retData := CreateKonsumensMessage(api)

	retDataEmail := EmailKonsumensMessage(api)

	if retData == "success" && retDataEmail == "success" {
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
