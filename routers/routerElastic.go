package routers

import (
	"kredit_plus/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/elastic", &controllers.ElasticController{}, "post:ElasticInsert")
	beego.Router("/elastic", &controllers.ElasticController{}, "get:SearchData")
	beego.Router("/elastic/all", &controllers.ElasticController{}, "get:ElasticGetAllData")
	beego.Router("/elastic/:id", &controllers.ElasticController{}, "delete:ElasticDelete")
	beego.Router("/elastic/:id", &controllers.ElasticController{}, "get:ElasticGetDataById")
}
