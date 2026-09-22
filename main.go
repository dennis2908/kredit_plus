package main

import (
	_ "fmt"
	loadconf "kredit_plus/LoadConf"
	_ "kredit_plus/routers"
	"kredit_plus/ssrf"
	"kredit_plus/token"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/lib/pq"
)

func init() {
	loadconf.Connects()
	beego.Debug("Filters init...")
	beego.InsertFilter("*", beego.BeforeRouter, token.Authenticate)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	orm.RegisterDriver("postgres", orm.DRPostgres)
	
	// Updated host from netwkreditPlus to localhost
	err := orm.RegisterDataBase("default",
		"postgres",
		"user=postgres password=123456 host=localhost port=5432 dbname=kredit_plus sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to register database: %v", err)
	}

	orm.RunSyncdb("default", false, true)
	orm.RunCommand()
}

func main() {
	err := ssrf.Main()
	if err != nil {
		log.Println(err)
		return
	}

	// Optimize Go runtime thread allocation
	numberOfCores := runtime.NumCPU()
	runtime.GOMAXPROCS(numberOfCores)

	// Initialize Cache
	_, err = cache.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache", "EmbedExpiry": "120"}`)
	if err != nil {
		logs.Error(err)
	}

	orm.Debug = true

	// Handle Port configuration
	log.Println("Env $PORT :", os.Getenv("PORT"))
	if envPort := os.Getenv("PORT"); envPort != "" {
		port, err := strconv.Atoi(envPort)
		if err != nil {
			log.Fatalf("Invalid $PORT value: %v", err)
		}
		log.Println("port : ", port)
		beego.BConfig.Listen.HTTPPort = port
		beego.BConfig.Listen.HTTPSPort = port
	}

	// Start web server (Blocking call)
	beego.Run()
}