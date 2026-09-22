package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"kredit_plus/models"
	"kredit_plus/token"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	beego.Controller
}

const (
	errCouldNotCreateToken = "could not create token"
	errUserNotFound        = "user not found"
	errInvalidToken        = "invalid token"
)

type userCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type userResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func toUserResponse(user *models.User) userResponse {
	return userResponse{Id: user.Id, Username: user.Username, Role: user.Role}
}

func (controller *AuthController) respond(status int, body interface{}) {
	controller.Ctx.Output.SetStatus(status)
	controller.Data["json"] = body
	controller.ServeJSON()
}

func (controller *AuthController) Register() {
	credentials := &userCredentials{}
	if err := json.Unmarshal(controller.Ctx.Input.RequestBody, credentials); err != nil || strings.TrimSpace(credentials.Username) == "" || credentials.Password == "" {
		controller.respond(http.StatusBadRequest, map[string]string{"error": "username and password are required"})
		return
	}

	username := strings.TrimSpace(credentials.Username)
	ormer := orm.NewOrm()
	var existing models.User
	if err := ormer.QueryTable("users").Filter("Username", username).One(&existing); err == nil {
		controller.respond(http.StatusConflict, map[string]string{"error": "username already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		controller.respond(http.StatusInternalServerError, map[string]string{"error": "could not create user"})
		return
	}
	user := &models.User{Username: username, PasswordHash: string(hash), Role: "user"}
	if _, err = ormer.Insert(user); err != nil {
		controller.respond(http.StatusInternalServerError, map[string]string{"error": "could not create user"})
		return
	}

	controller.respond(http.StatusCreated, toUserResponse(user))
}

func (controller *AuthController) Login() {
	credentials := &userCredentials{}
	if err := json.Unmarshal(controller.Ctx.Input.RequestBody, credentials); err != nil {
		controller.respond(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	var user models.User
	ormer := orm.NewOrm()
	if err := ormer.QueryTable("users").Filter("Username", strings.TrimSpace(credentials.Username)).One(&user); err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)) != nil {
		controller.respond(http.StatusUnauthorized, map[string]string{"error": "invalid username or password"})
		return
	}

	accessToken, err := token.Create(user.Id, user.Username, user.Role)
	if err != nil {
		controller.respond(http.StatusInternalServerError, map[string]string{"error": errCouldNotCreateToken})
		return
	}
	refreshToken, err := token.CreateRefresh(user.Id, user.Username, user.Role)
	if err != nil {
		controller.respond(http.StatusInternalServerError, map[string]string{"error": errCouldNotCreateToken})
		return
	}
	controller.respond(http.StatusOK, map[string]interface{}{"token": accessToken, "refresh_token": refreshToken, "user": toUserResponse(&user)})
}

func (controller *AuthController) Refresh() {
	request := &refreshRequest{}
	if err := json.Unmarshal(controller.Ctx.Input.RequestBody, request); err != nil || request.RefreshToken == "" {
		controller.respond(http.StatusBadRequest, map[string]string{"error": "refresh_token is required"})
		return
	}
	userID, ok := token.UserIDFromToken(request.RefreshToken, "refresh")
	if !ok {
		controller.respond(http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}
	var user models.User
	if err := orm.NewOrm().QueryTable("users").Filter("Id", userID).One(&user); err != nil {
		controller.respond(http.StatusUnauthorized, map[string]string{"error": errUserNotFound})
		return
	}
	accessToken, err := token.Create(user.Id, user.Username, user.Role)
	if err != nil {
		controller.respond(http.StatusInternalServerError, map[string]string{"error": errCouldNotCreateToken})
		return
	}
	controller.respond(http.StatusOK, map[string]string{"token": accessToken})
}

func (controller *AuthController) Me() {
	userID, ok := token.UserID(controller.Ctx)
	if !ok {
		controller.respond(http.StatusUnauthorized, map[string]string{"error": errInvalidToken})
		return
	}

	var user models.User
	if err := orm.NewOrm().QueryTable("users").Filter("Id", userID).One(&user); err != nil {
		controller.respond(http.StatusNotFound, map[string]string{"error": errUserNotFound})
		return
	}
	controller.respond(http.StatusOK, toUserResponse(&user))
}

func (controller *AuthController) Update() {
	userID, ok := token.UserID(controller.Ctx)
	if !ok {
		controller.respond(http.StatusUnauthorized, map[string]string{"error": errInvalidToken})
		return
	}
	credentials := &userCredentials{}
	if err := json.Unmarshal(controller.Ctx.Input.RequestBody, credentials); err != nil {
		controller.respond(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	ormer := orm.NewOrm()
	var user models.User
	if err := ormer.QueryTable("users").Filter("Id", userID).One(&user); err != nil {
		controller.respond(http.StatusNotFound, map[string]string{"error": errUserNotFound})
		return
	}
	if strings.TrimSpace(credentials.Username) != "" {
		user.Username = strings.TrimSpace(credentials.Username)
	}
	if credentials.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
		if err != nil {
			controller.respond(http.StatusInternalServerError, map[string]string{"error": "could not update user"})
			return
		}
		user.PasswordHash = string(hash)
	}
	if _, err := ormer.Update(&user, "Username", "PasswordHash"); err != nil {
		controller.respond(http.StatusConflict, map[string]string{"error": "could not update user"})
		return
	}
	controller.respond(http.StatusOK, toUserResponse(&user))
}

func (controller *AuthController) Delete() {
	userID, ok := token.UserID(controller.Ctx)
	if !ok {
		controller.respond(http.StatusUnauthorized, map[string]string{"error": errInvalidToken})
		return
	}
	if _, err := orm.NewOrm().QueryTable("users").Filter("Id", userID).Delete(); err != nil {
		controller.respond(http.StatusInternalServerError, map[string]string{"error": "could not delete user"})
		return
	}
	controller.Ctx.Output.SetStatus(http.StatusNoContent)
}
