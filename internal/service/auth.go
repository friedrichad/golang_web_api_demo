package service

import (
	"encoding/base64"
	"log"
	"strings"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Authentication(c *gin.Context) (*model.TokenResponse, *common.Error)
}

type AuthService struct {
	repository          repository.IUserRepository
	accessTokenExpired  int
	refreshTokenExpired int
	clientId            string
	basicAuth           string
	jwtSecret           string
	clientTypeNoExp     []string
}

func NewAuthService() IAuthService {
	return &AuthService{
		repository:          repository.NewUserRepository(),
		accessTokenExpired:  viper.GetInt("oauth.access-token-expired"),
		refreshTokenExpired: viper.GetInt("oauth.refresh-token-expired"),
		clientId:            viper.GetString("oauth.client-id"),
		basicAuth:           base64.StdEncoding.EncodeToString([]byte(viper.GetString("oauth.client-id") + ":" + viper.GetString("oauth.client-secret"))),
		jwtSecret:           viper.GetString("oauth.jwt-secret"),
		clientTypeNoExp:     viper.GetStringSlice("oauth.client_type_no_expire"),
	}
}

func (a AuthService) Authentication(c *gin.Context) (*model.TokenResponse, *common.Error) {
	grantType := c.Request.FormValue("grant_type")
	if grantType == "password" {
		return createNewToken(c, a)
	}
	if grantType == "refresh_token" {
		return refreshToken(c, a)
	}
	return nil, &common.Error{Code: "400", Message: "Grant type không hỗ trợ"}
}

func createNewToken(c *gin.Context, a AuthService) (*model.TokenResponse, *common.Error) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	if len(password) == 0 {
		return nil, common.AuthenticationFail
	}
	user, err := a.repository.GetByUsername(username)
	if err != nil {
		log.Printf("Error when get user by username: %s", err.Error())
		return nil, common.AuthenticationFail
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, common.AuthenticationFail
	}
	response := &model.TokenResponse{
		AccessToken: "",
		TokenType:   "bearer",
	}
	response.Id = user.UserUUID
	response.Username = user.Username
	response.Active = true
	response.Exp = getExpiredTime(a.accessTokenExpired)
	response.RefreshExp = getExpiredTime(a.refreshTokenExpired)
	authorities, err := a.repository.GetAuthorities(user.UserUUID)
	if err != nil && err.Error() != "record not found" {
		return nil, common.SystemError
	}
	response.Authorities = authorities
	accessToken, err := createJwtToken(a.jwtSecret, *response)
	if err != nil {
		return nil, common.SystemError
	}
	response.AccessToken = accessToken
	return response, nil
}

func getExpiredTime(expTime int) int64 {
	expAccessToken := time.Duration(expTime * int(time.Second))
	return (time.Now().UnixNano() + int64(expAccessToken)) / int64(time.Second)
}

func createJwtToken(jwtSecret string, token model.TokenResponse) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)
	claims := t.Claims.(jwt.MapClaims)
	claims["id"] = token.Id
	claims["user_name"] = token.Username
	claims["exp"] = token.Exp
	claims["authorities"] = token.Authorities
	claims["client_id"] = token.ClientId
	claims["refresh_exp"] = token.RefreshExp
	accessToken, err := t.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func refreshToken(c *gin.Context, a AuthService) (*model.TokenResponse, *common.Error) {
	claims, ok := extractClaims(c.Request.FormValue("refresh_token"), []byte(a.jwtSecret))
	if !ok {
		return nil, common.TokenInvalid
	}
	response := &model.TokenResponse{
		AccessToken: "",
		TokenType:   "bearer",
	}
	response.Id = claims.Id
	response.Username = claims.Username
	response.Active = true
	response.Exp = getExpiredTime(a.accessTokenExpired)
	response.RefreshExp = getExpiredTime(a.refreshTokenExpired)
	response.Authorities = claims.Authorities
	accessToken, err := createJwtToken(a.jwtSecret, *response)
	if err != nil {
		return nil, common.SystemError
	}
	response.AccessToken = accessToken
	return response, nil
}

func extractClaims(tokenStr string, hmacSecret []byte) (*model.Claims, bool) {
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)
	claims := &model.Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err != nil && err.Error() != "Token is expired" {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
	return claims, claims.RefreshTokenExpired()
}
