package main

import (
	_ "fmt"
	Connection "kredit_plus/connects"
	_ "kredit_plus/routers"
	"kredit_plus/ssrf"
	"kredit_plus/token"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/lib/pq"
)

func init() { // init instead of int

	Connection.Connects()
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

		beego.InsertFilter("/konsumen/*", beego.BeforeExec, token.ValidateToken)
		beego.InsertFilter("/transaksi/*", beego.BeforeExec, token.ValidateToken)
		beego.InsertFilter("/transaksidetails/*", beego.BeforeExec, token.ValidateToken)
		beego.InsertFilter("/konsumen_mongo_insert/*", beego.BeforeExec, token.ValidateToken)
		beego.InsertFilter("/konsumen_mongo_update/*", beego.BeforeExec, token.ValidateToken)
		beego.InsertFilter("/konsumen_excel_read/*", beego.BeforeExec, token.ValidateToken)

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
