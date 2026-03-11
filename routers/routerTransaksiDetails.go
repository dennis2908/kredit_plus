package routers

import (
	"kredit_plus/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/transaksidetails/:id", &controllers.TransaksiDetailsController{}, "put:UpdateTransaksiDetails")
	beego.Router("/transaksidetails/", &controllers.TransaksiDetailsController{}, "post:CreateTransaksiDetails")
	beego.Router("/transaksidetails/", &controllers.TransaksiDetailsController{}, "get:GetAllTransaksiDetails")
	beego.Router("/transaksidetails/:id", &controllers.TransaksiDetailsController{}, "get:GetTransaksiDetailsByIdKonsumen")

}
