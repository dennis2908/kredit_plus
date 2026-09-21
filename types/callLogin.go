package types

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Login interface {
	CekLogin() (string, string)
}

type DoLogin struct {
	Calllogin Login
}

type MakeLogin struct{}

func (makeLogin MakeLogin) CekLogin() (string, string) {
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	token, _ := tokenString.SignedString([]byte(os.Getenv("token_secret_key")))

	tokenRefesh, _ := tokenString.SignedString([]byte(os.Getenv("token_refresh_key")))

	vs := base64.URLEncoding.EncodeToString([]byte(tokenRefesh))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(os.Getenv(("cookie_secret_key"))))
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")

	return token, cookie

}
