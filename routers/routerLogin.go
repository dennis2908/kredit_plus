package routers

import (
	"kredit_plus/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/refresh/token", &controllers.AuthController{}, "post:Refresh")
	beego.Router("/users", &controllers.AuthController{}, "post:Register")
	beego.Router("/users/me", &controllers.AuthController{}, "get:Me;put:Update;delete:Delete")
}
