package connects

import (
	_ "fmt"
	_ "kredit_plus/routers"

	loadconf "kredit_plus/LoadConf"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/lib/pq"
)

func Connects() { // init instead of int

	loadconf.Connects()

	// rabbits.GetData()
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
