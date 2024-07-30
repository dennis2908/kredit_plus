package token

import (
	"fmt"

	"github.com/astaxie/beego/context"

	"github.com/dgrijalva/jwt-go"

	"strings"
)

var secretKey = []byte("secret-key")

// JwtAuth 中间件，检查token
func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func ValidateToken(ctx *context.Context) {
	tokenString := ctx.Input.Header("Authorization")

	if tokenString == "" {
		ctx.Output.SetStatus(401)
		ctx.ResponseWriter.Write([]byte("need auth"))
		return
	}
	err := verifyToken(strings.Split(tokenString, " ")[1])
	if err != nil {
		ctx.Output.SetStatus(401)
		ctx.ResponseWriter.Write([]byte("need auth"))
		return
	}

}