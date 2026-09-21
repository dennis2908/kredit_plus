package routers

import (
	"kredit_plus/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/users", &controllers.AuthController{}, "post:CreateUser")
	beego.Router("/users", &controllers.AuthController{}, "get:GetAllUsers")
	beego.Router("/users/:id", &controllers.AuthController{}, "put:UpdateUser")
	beego.Router("/users/:id", &controllers.AuthController{}, "get:GetUserByID")
}
