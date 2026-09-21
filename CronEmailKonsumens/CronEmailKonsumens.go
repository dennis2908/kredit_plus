package main

import (
	"fmt"
	_ "fmt"
	Connection "kredit_plus/connects"
	"kredit_plus/models"
	_ "kredit_plus/routers"
	"log"
	"net/smtp"

	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"

	"os"
	"os/signal"
	"syscall"
	"time"

	cron "github.com/robfig/cron/v3"
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
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakartaTime))

	// stop scheduler tepat sebelum fungsi berakhir
	defer scheduler.Stop()

	// set task yang akan dijalankan scheduler
	// gunakan crontab string untuk notifikasi harian tiap jam 9 malam ke email masing masing konsumen
	scheduler.AddFunc("0 21 * * 1-7", NotifyDailyNightNotif)

	// scheduler.AddFunc("*/1 * * * *", NotifyDailyNightNotif)
	// start scheduler
	go scheduler.Start()

	// trap SIGINT untuk trigger shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func Email(Full_name string, Email string) {

	if Full_name == "" {
		Full_name = "Konsumen"
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	// Subject and body
	subject := "Good Night, " + Full_name
	body := "Hi, Good Night " + Full_name

	// Sender and receiver
	from := os.Getenv("SMTP_FROM")
	to := []string{
		Email,
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

func NotifyDailyNightNotif() {
	o := orm.NewOrm()
	o.Using("default")
	var konsumens []models.Konsumens
	_, err := o.QueryTable("konsumens").All(&konsumens)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, data := range konsumens {
		if data.Email != "" {
			Email(data.Full_name, data.Email)
		}

	}
}
