package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
)

func Create(userID int, username, role string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      strconv.Itoa(userID),
		"username": username,
		"role":     role,
		"kind":     "access",
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}).SignedString([]byte(os.Getenv("token_secret_key")))
}

func CreateRefresh(userID int, username, role string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      strconv.Itoa(userID),
		"username": username,
		"role":     role,
		"kind":     "refresh",
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}).SignedString([]byte(os.Getenv("token_secret_key")))
}

func parse(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("token_secret_key")
	if secret == "" {
		return nil, fmt.Errorf("token secret is not configured")
	}
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
}

func Authenticate(ctx *context.Context) {
	if ctx.Request.Method == "OPTIONS" || ctx.Input.URL() == "/login" || ctx.Input.URL() == "/users" || ctx.Input.URL() == "/refresh/token" {
		return
	}
	authorization := strings.Fields(ctx.Input.Header("Authorization"))
	if len(authorization) != 2 || !strings.EqualFold(authorization[0], "Bearer") {
		ctx.Output.SetStatus(401)
		ctx.ResponseWriter.Write([]byte(`{"error":"authentication required"}`))
		return
	}
	parsed, err := parse(authorization[1])
	if err != nil {
		if strings.Contains(err.Error(), "token secret is not configured") {
			ctx.Output.SetStatus(500)
			ctx.ResponseWriter.Write([]byte(`{"error":"auth configuration error: token secret is not configured"}`))
			return
		}
		ctx.Output.SetStatus(401)
		ctx.ResponseWriter.Write([]byte(`{"error":"invalid token"}`))
		return
	}
	if !parsed.Valid {
		ctx.Output.SetStatus(401)
		ctx.ResponseWriter.Write([]byte(`{"error":"invalid token"}`))
	}
}

func ValidateToken(ctx *context.Context) {
	Authenticate(ctx)
}

func UserID(ctx *context.Context) (int, bool) {
	authorization := strings.Fields(ctx.Input.Header("Authorization"))
	if len(authorization) != 2 || !strings.EqualFold(authorization[0], "Bearer") {
		return 0, false
	}
	return UserIDFromToken(authorization[1], "access")
}

func UserIDFromToken(tokenString, requiredKind string) (int, bool) {
	parsed, err := parse(tokenString)
	if err != nil || !parsed.Valid {
		return 0, false
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}
	if kind, ok := claims["kind"].(string); !ok || kind != requiredKind {
		return 0, false
	}
	value, ok := claims["sub"].(string)
	if !ok {
		return 0, false
	}
	id, err := strconv.Atoi(value)
	return id, err == nil && id > 0
}
