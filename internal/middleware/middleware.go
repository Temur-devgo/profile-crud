package middlware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"sarkor/internal/handler"
	"sarkor/internal/service"
)

func ValidateToken(ctx *gin.Context) {
	token, err := ctx.Cookie("token")
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims := &service.CustomClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("my_secret_key"), nil
	})

	if err != nil || err == jwt.ErrSignatureInvalid {

		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !tkn.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "user", handler.UserTokenData{UserID: claims.UserID, Login: claims.Login}))

	ctx.Next()
}
