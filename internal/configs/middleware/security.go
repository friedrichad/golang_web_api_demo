package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

func BasicAuthenticator() gin.HandlerFunc {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(viper.GetString("oauth.client-id") + ":" + viper.GetString("oauth.client-secret")))
	return func(c *gin.Context) {
		basic := c.Request.Header["Authorization"]
		if len(basic) == 0 || strings.Replace(basic[0], "Basic ", "", 1) != basicAuth {
			c.JSON(http.StatusUnauthorized, common.Error{Code: "401", Message: "Client không hợp lệ"})
			c.Abort()
		}
	}
}

func BearerAuthenticator() gin.HandlerFunc {
	hmacSecret := []byte(viper.GetString("oauth.jwt-secret"))
	return func(c *gin.Context) {
		jwtClaims, valid := extractClaims(c.GetHeader("Authorization"), hmacSecret)
		if !valid {
			c.JSON(http.StatusUnauthorized, common.Error{Code: "401", Message: "Token không hợp lệ hoặc đã hết hạn"})
			c.Abort()
			return
		}
		c.Set("user_id", jwtClaims.Id)
		c.Set("username", jwtClaims.Username)
		c.Set("authorities", jwtClaims.Authorities)
	}
}

func Authorizator(authority ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorities := c.GetStringSlice("authorities")
		if len(authorities) == 0 || !utils.AnyContains(authorities, authority) {
			c.JSON(http.StatusForbidden, common.Error{Code: "403", Message: "Không có quyền truy cập"})
			c.Abort()
			return
		}
	}
}

func extractClaims(tokenStr string, hmacSecret []byte) (*model.Claims, bool) {
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)
	claims := &model.Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err == nil {
		return claims, true
	}
	log.Printf("Invalid JWT Token")
	return nil, false
}
