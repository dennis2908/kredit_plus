package controllers

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "fmt"
	models "kredit_plus/models"
	mongoconn "kredit_plus/mongoconn"
	redisconn "kredit_plus/redisconn"
	"kredit_plus/structs"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	_ "github.com/leekchan/accounting"
	amqp "github.com/rabbitmq/amqp091-go"
	_ "github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type KonsumensController struct {
	beego.Controller
}

func (api *KonsumensController) GetAllKonsumens() {
	o := orm.NewOrm()
	o.Using("default")
	var konsumens []models.Konsumens
	num, err := o.QueryTable("konsumens").All(&konsumens)
	rdb, _ := redisconn.Connect()
	urlsJson, _ := json.Marshal(konsumens)
	token, _ := GenerateRandomString(32)

	ttl := time.Duration(3) * time.Second

	op1 := rdb.Set(context.Background(), token, urlsJson, ttl)
	if err := op1.Err(); err != nil {
		fmt.Printf("unable to SET data. error: %v", err)
		return
	}
	op2 := rdb.Get(context.Background(), token)
	fmt.Printf("data", op2)

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

func (api *KonsumensController) GetRefreshToken() {
	val := api.Ctx.GetCookie("token")

	var strErr = "refresh token not available"
	if val == "" {
		api.Data["json"] = strErr
	} else {

		parts := strings.SplitN(val, "|", 3)

		if len(parts) != 3 {
			api.Data["json"] = strErr
		} else {

			vs := parts[0]
			timestamp := parts[1]
			sig := parts[2]

			h := hmac.New(sha1.New, []byte(os.Getenv(("cookie_secret_key"))))
			fmt.Fprintf(h, "%s%s", vs, timestamp)

			if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
				api.Data["json"] = strErr
			} else {

				res, _ := base64.URLEncoding.DecodeString(vs)

				token, err := jwt.Parse(string(res), func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("token_refresh_key")), nil
				})

				if err != nil {
					api.Data["json"] = strErr
				}

				if !token.Valid {
					api.Data["json"] = strErr
				}

				if err == nil {

					tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
						"exp": time.Now().Add(time.Hour * 24).Unix(),
					})
					token, _ := tokenString.SignedString([]byte(os.Getenv("token_secret_key")))

					api.Data["json"] = "{'your new token key': '" + token + "'}"
				}
			}
		}
	}

	api.ServeJSON()
}

func (api *KonsumensController) GetToken() {
	login := &models.Login{}
	json.Unmarshal(api.Ctx.Input.RequestBody, login)

	fmt.Println(1111, login.Username)
	fmt.Println(2222, login.Password)

	api.Data["json"] = "wrong username or password"
	if login.Username == "admin" && login.Password == "mypassword" {
		tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		token, _ := tokenString.SignedString([]byte(os.Getenv("token_secret_key")))

		tokenRefesh, _ := tokenString.SignedString([]byte(os.Getenv("token_refresh_key")))

		vs := base64.URLEncoding.EncodeToString([]byte(tokenRefesh))
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		h := hmac.New(sha1.New, []byte(os.Getenv(("cookie_secret_key"))))
		fmt.Fprintf(h, "%s%s", vs, timestamp)
		fmt.Println(tokenRefesh)
		fmt.Printf("mystr:\t %v \n", []byte(os.Getenv("token_secret_key")))
		fmt.Printf("mystr:\t %v \n", []byte(os.Getenv("token_refresh_key")))
		sig := fmt.Sprintf("%02x", h.Sum(nil))
		cookie := strings.Join([]string{vs, timestamp, sig}, "|")

		api.Ctx.SetCookie("token", cookie, time.Now().Add(time.Hour*24).Unix(), "/")
		api.Data["json"] = "{'tokenString': '" + token + "'}"
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

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func (api *KonsumensController) GetKonsumensById() {
	o := orm.NewOrm()
	o.Using("default")
	// ctx := context.Background()
	var konsumens []models.Konsumens
	idInt, _ := strconv.Atoi(api.Ctx.Input.Param(":id"))
	num, err := o.QueryTable("konsumens").Filter("id", idInt).All(&konsumens)
	rdb, _ := redisconn.Connect()
	urlsJson, _ := json.Marshal(konsumens)
	token, _ := GenerateRandomString(32)

	ttl := time.Duration(3) * time.Second

	op1 := rdb.Set(context.Background(), token, urlsJson, ttl)
	if err := op1.Err(); err != nil {
		fmt.Printf("unable to SET data. error: %v", err)
		return
	}
	rdb.Get(context.Background(), token)
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

func InsertKonsumensExcelMessage(modelKonsumen models.Konsumens) string {
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
	msg := failOnError(err, "Failed to connect to RabbitMQ")
	if msg == "Error" {
		return "Error"
	}
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"insertExcelKonsumens", // name
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

	dataSend, err := json.Marshal(modelKonsumen)
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

func UpdateKonsumensMessage(api *KonsumensController) string {
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
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
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
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
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
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

func excelDateToDate(excelDate string) time.Time {
	in, _ := strconv.ParseFloat(strings.TrimSpace(excelDate), 64)
	excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	tm := excelEpoch.Add(time.Duration(in * float64(24*time.Hour)))
	return tm
}

func (api *KonsumensController) ReadExcelKonsumens() {

	xlsx, err := excelize.OpenFile("static/excel/" + api.Ctx.Input.Param(":filename") + ".xlsx")
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	rowsExcel, _ := xlsx.GetRows("Sheet1")

	rows := make([]models.Konsumens, 0)
	for i, rowsExcel := range rowsExcel {
		if i == 0 {
			// Skip header row
			continue
		}
		rowEmail := rowsExcel[0]
		rowNik := rowsExcel[1]
		rowFullName := rowsExcel[2]
		rowLegalName := rowsExcel[3]
		rowPlace_birth := rowsExcel[4]
		// rowDateBirth, _ := xlsx.GetCellValue(sheet1Name, fmt.Sprintf("F%d", i))
		// rowDateBirthX := excelDateToDate(rowsExcel[5])
		rowDateBirthX, _ := strconv.ParseFloat(rowsExcel[5], 64)
		// rowDateBirth := excelDateToDate(rowsExcel[5])
		rowSalary, _ := strconv.Atoi(rowsExcel[6])
		rowFotoKTP := rowsExcel[7]
		rowFotoSelfie := rowsExcel[8]
		rowDateBirth, _ := time.Parse("2006-01-02", rowsExcel[5])

		log.Printf(" [x] Sent %s\n", rowDateBirthX)

		InsertKonsumensExcelMessage(models.Konsumens{Konsumen: models.Konsumen{Nik: rowNik, Full_name: rowFullName, Email: rowEmail, Legal_name: rowLegalName, Place_birth: rowPlace_birth, Salary: rowSalary, Foto_ktp: rowFotoKTP, Foto_selfie: rowFotoSelfie}, Date_birth: rowDateBirth})
	}

	fmt.Printf("%v \n", rows)
	// if AllKonsumensCheck(api) != "" {
	// 	api.Data["json"] = AllKonsumensCheck(api)
	// 	api.ServeJSON()
	// 	return
	// }

	// retData := CreateKonsumensMessage(api)

	// retDataEmail := EmailKonsumensMessage(api)

	// if retData == "success" && retDataEmail == "success" {
	// 	api.Data["json"] = "Successfully insert data"
	// 	api.Ctx.ResponseWriter.WriteHeader(200)
	// 	api.ServeJSON()
	// } else {
	// 	api.Ctx.ResponseWriter.WriteHeader(500)
	// 	api.Data["json"] = "Error"

	// 	api.ServeJSON()
	// }
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
