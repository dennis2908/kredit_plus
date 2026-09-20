package main

import (
	_ "fmt"
	_ "kredit_plus/routers"
	"kredit_plus/ssrf"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/lib/pq"
)

func init() { // init instead of int
	beego.Debug("Filters init...")

	// CORS for https://foo.* origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default",
		"postgres",
		"user=postgres password=123456 host=127.0.0.1 port=5432 dbname=kredit_plus sslmode=disable")
	orm.RunSyncdb("default", false, true)
	orm.RunCommand()
}
func main() {
	err := ssrf.Main()
	if err != nil {
		log.Println(err)
		return
	}
	numberOfCores := runtime.NumCPU()
	runtime.GOMAXPROCS(numberOfCores)
	var wg sync.WaitGroup
	for i := 0; i < numberOfCores; i++ {
		wg.Add(1)
		_, err := cache.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache", "EmbedExpiry": "120"}`)

		orm.Debug = true

		o := orm.NewOrm()
		o.Using("default")

		if err != nil {
			logs.Error(err)
		}
		log.Println("Env $PORT :", os.Getenv("PORT"))
		if os.Getenv("PORT") != "" {
			port, err := strconv.Atoi(os.Getenv("PORT"))
			if err != nil {
				log.Fatal(err)
				log.Fatal("$PORT must be set")
			}
			log.Println("port : ", port)
			beego.BConfig.Listen.HTTPPort = port
			beego.BConfig.Listen.HTTPSPort = port
		}

		beego.Run()
	}
	wg.Wait()

}
