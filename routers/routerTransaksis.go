package routers

import (
	"kredit_plus/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/transaksi/:id", &controllers.TransaksisController{}, "put:UpdateTransaksis")
	beego.Router("/transaksi/", &controllers.TransaksisController{}, "post:CreateTransaksis")
	beego.Router("/transaksi/", &controllers.TransaksisController{}, "get:GetAllTransaksis")
	beego.Router("/transaksi/:id", &controllers.TransaksisController{}, "get:GetTransaksisById")
	beego.Router("/transaksi/konsumen/:id", &controllers.TransaksisController{}, "get:GetTransaksiByIdKons")
	
}
