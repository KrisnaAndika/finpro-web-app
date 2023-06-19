package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		// TODO: answer here
		cookie, err := ctx.Request.Cookie("session_token")
		if err != nil {
			if ctx.GetHeader("Content-Type") == "application/json" {
				ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(err.Error()))
			} else {
				ctx.Redirect(http.StatusSeeOther, "/login")
			}
			return
		}

		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(tkn *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
			return
		}

		if token.Valid == false {
			ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse("token tidak valid"))
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Next()
	})
}
