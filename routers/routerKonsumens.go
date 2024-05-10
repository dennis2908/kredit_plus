package routers

import (
	"kredit_plus/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/konsumen/:id", &controllers.KonsumensController{}, "put:UpdateKonsumens")
	beego.Router("/konsumen/", &controllers.KonsumensController{}, "post:CreateKonsumens")
	beego.Router("/konsumen_excel/:filename", &controllers.KonsumensController{}, "get:ExcelKonsumens")
	beego.Router("/konsumen/", &controllers.KonsumensController{}, "get:GetAllKonsumens")
	beego.Router("/konsumen_mongo_insert/", &controllers.KonsumensController{}, "get:GetAllKonsumensMongoInsert")
	beego.Router("/konsumen_mongo_update/", &controllers.KonsumensController{}, "get:GetAllKonsumensMongoUpdate")
	beego.Router("/konsumen/:id", &controllers.KonsumensController{}, "get:GetKonsumensById")

}
