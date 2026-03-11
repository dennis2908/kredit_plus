package routers

import (
	"kredit_plus/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controllers.KonsumensController{}, "post:GetToken")
	beego.Router("/refresh/token", &controllers.KonsumensController{}, "post:GetRefreshToken")

}
