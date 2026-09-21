package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	models "kredit_plus/models"
	"kredit_plus/structs"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const accessTokenLifetime = 24 * time.Hour

type AuthController struct {
	beego.Controller
}

func (api *AuthController) CreateUser() {
	var request structs.UserRequest
	if !decodeUserRequest(api, &request, true) {
		return
	}

	o := orm.NewOrm()
	user := models.User{
		Email:    normalizeEmail(request.Email),
		FullName: strings.TrimSpace(request.FullName),
		IsActive: true,
	}
	if request.IsActive != nil {
		user.IsActive = *request.IsActive
	}
	var err error
	user.PasswordHash, err = HashPassword(request.Password)
	if err != nil {
		writeJSONError(api, http.StatusInternalServerError, "could not hash password")
		return
	}
	if _, err = o.Insert(&user); err != nil {
		writeJSONError(api, http.StatusConflict, "email already exists")
		return
	}

	api.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	api.Data["json"] = toPublicUser(user)
	api.ServeJSON()
}

func (api *AuthController) UpdateUser() {
	userID, err := strconv.Atoi(api.Ctx.Input.Param(":id"))
	if err != nil || userID < 1 {
		writeJSONError(api, http.StatusBadRequest, "invalid user id")
		return
	}
	var request structs.UserRequest
	if !decodeUserRequest(api, &request, false) {
		return
	}

	o := orm.NewOrm()
	var user models.User
	if err = o.QueryTable(new(models.User)).Filter("id", userID).One(&user); err != nil {
		writeJSONError(api, http.StatusNotFound, "user not found")
		return
	}
	if request.Email != "" {
		user.Email = normalizeEmail(request.Email)
	}
	if request.FullName != "" {
		user.FullName = strings.TrimSpace(request.FullName)
	}
	if request.IsActive != nil {
		user.IsActive = *request.IsActive
	}
	if request.Password != "" {
		user.PasswordHash, err = HashPassword(request.Password)
		if err != nil {
			writeJSONError(api, http.StatusInternalServerError, "could not hash password")
			return
		}
	}
	if _, err = o.Update(&user); err != nil {
		writeJSONError(api, http.StatusConflict, "email already exists")
		return
	}

	api.Data["json"] = toPublicUser(user)
	api.ServeJSON()
}

func (api *AuthController) GetAllUsers() {
	var users []models.User
	if _, err := orm.NewOrm().QueryTable(new(models.User)).All(&users); err != nil {
		writeJSONError(api, http.StatusInternalServerError, "could not load users")
		return
	}
	result := make([]structs.UserPublic, 0, len(users))
	for _, user := range users {
		result = append(result, toPublicUser(user))
	}
	api.Data["json"] = result
	api.ServeJSON()
}

func (api *AuthController) GetUserByID() {
	userID, err := strconv.Atoi(api.Ctx.Input.Param(":id"))
	if err != nil || userID < 1 {
		writeJSONError(api, http.StatusBadRequest, "invalid user id")
		return
	}
	var user models.User
	if err = orm.NewOrm().QueryTable(new(models.User)).Filter("id", userID).One(&user); err != nil {
		writeJSONError(api, http.StatusNotFound, "user not found")
		return
	}
	api.Data["json"] = toPublicUser(user)
	api.ServeJSON()
}

func (api *AuthController) Login() {
	var request structs.LoginRequest
	if err := json.Unmarshal(api.Ctx.Input.RequestBody, &request); err != nil || strings.TrimSpace(request.Email) == "" || request.Password == "" {
		api.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		api.Data["json"] = map[string]string{"error": "email and password are required"}
		api.ServeJSON()
		return
	}

	var user models.User
	o := orm.NewOrm()
	err := o.QueryTable(new(models.User)).Filter("email", strings.ToLower(strings.TrimSpace(request.Email))).Filter("is_active", true).One(&user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil {
		api.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		api.Data["json"] = map[string]string{"error": "invalid email or password"}
		api.ServeJSON()
		return
	}

	expiresAt := time.Now().Add(accessTokenLifetime)
	claims := jwt.MapClaims{
		"sub":   user.Id,
		"email": user.Email,
		"exp":   expiresAt.Unix(),
		"iat":   time.Now().Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret())
	if err != nil {
		api.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		api.Data["json"] = map[string]string{"error": "could not create access token"}
		api.ServeJSON()
		return
	}

	api.Data["json"] = structs.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(accessTokenLifetime.Seconds()),
		User: structs.UserPublic{
			Id:       user.Id,
			Email:    user.Email,
			FullName: user.FullName,
		},
	}
	api.ServeJSON()
}

func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "development-only-change-this-secret"
	}
	return []byte(secret)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func ValidateBearerToken(authorization string) (*jwt.Token, error) {
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return nil, jwt.ErrSignatureInvalid
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret(), nil
	})
}

func decodeUserRequest(api *AuthController, request *structs.UserRequest, passwordRequired bool) bool {
	if err := json.Unmarshal(api.Ctx.Input.RequestBody, request); err != nil {
		writeJSONError(api, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	if strings.TrimSpace(request.Email) == "" || (passwordRequired && request.Password == "") {
		writeJSONError(api, http.StatusBadRequest, "email and password are required")
		return false
	}
	if request.Email != "" && !strings.Contains(request.Email, "@") {
		writeJSONError(api, http.StatusBadRequest, "invalid email")
		return false
	}
	return true
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func toPublicUser(user models.User) structs.UserPublic {
	return structs.UserPublic{Id: user.Id, Email: user.Email, FullName: user.FullName}
}

func writeJSONError(api *AuthController, status int, message string) {
	api.Ctx.ResponseWriter.WriteHeader(status)
	api.Data["json"] = map[string]string{"error": message}
	api.ServeJSON()
}
