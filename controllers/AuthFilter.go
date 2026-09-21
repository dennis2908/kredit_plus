package controllers

import (
	"net/http"

	"github.com/astaxie/beego/context"
)

func RequireAuthentication(ctx *context.Context) {
	if ctx.Input.Method() == http.MethodOptions || isPublicAuthPath(ctx.Input.URL()) {
		return
	}

	token, err := ValidateBearerToken(ctx.Input.Header("Authorization"))
	if err == nil && token.Valid {
		return
	}

	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.Abort(http.StatusUnauthorized, `{"error":"authentication required"}`)
}

func isPublicAuthPath(path string) bool {
	return path == "/login"
}
