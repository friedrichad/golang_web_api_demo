package service

import (
	"encoding/base64"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/redis"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Authentication(c *gin.Context) (*model.TokenResponse, *common.Error)
	Register(c *gin.Context) (*dtos.UserResponse, *common.Error)
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
	response.Id = strconv.FormatInt(int64(user.UserID), 10)
	response.Username = user.Username
	response.Active = true
	response.Exp = getExpiredTime(a.accessTokenExpired)
	response.RefreshExp = getExpiredTime(a.refreshTokenExpired)
	authorities, err := a.repository.GetAuthorities(user.UserID)
	if err != nil && err.Error() != "record not found" {
		return nil, common.SystemError
	}
	response.Authorities = authorities
	accessToken, err := createJwtToken(a.jwtSecret, *response)
	if err != nil {
		return nil, common.SystemError
	}
	response.AccessToken = accessToken
	ttl := time.Until(time.Unix(response.Exp, 0))
	err = redis.Save(redis.Rdb, "auth:token:"+response.Id, response.AccessToken, ttl)
	if err != nil {
		log.Printf("Không lưu được token vào Redis: %v", err)
		return nil, common.SystemError
	}

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
	// Save new token to Redis with TTL
	ttl := time.Until(time.Unix(response.RefreshExp, 0))
	err = redis.Save(redis.Rdb, "auth:token:"+response.Id, response.AccessToken, ttl)
	if err != nil {
		log.Printf("Không lưu được token vào Redis: %v", err)
		return nil, common.SystemError
	}
	return response, nil
}

func extractClaims(tokenStr string, hmacSecret []byte) (*model.Claims, bool) {
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)
	claims := &model.Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err != nil && err.Error() != "Token is expired" {
		log.Printf("Invalid JWT Token: %v", err)
		return nil, false
	}
	// For refresh token: Allow expired access token, but check if refresh_exp is still valid
	if !claims.IsRefreshTokenValid() {
		log.Printf("Refresh token has expired")
		return nil, false
	}
	stored, err := redis.Exists(redis.Rdb, "auth:token:" +claims.Id)
	if err != nil {
    	log.Printf("Lỗi khi check Redis: %v", err)
    	return nil, false
	}
	if !stored {
    	log.Printf("Refresh token không tồn tại trong Redis (có thể đã bị revoke)")
    	return nil, false
	}
    return claims, true
}

func (a AuthService) Register(c *gin.Context) (*dtos.UserResponse, *common.Error) {
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	existingUser, _ := a.repository.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, &common.Error{Code: "400", Message: "Username đã tồn tại"}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, common.SystemError
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, common.SystemError
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	userRepoTx := a.repository.(*repository.UserRepository).WithTx(tx)

	user := &model.User{
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		StatusInt:    1,
		CreatedAt:    time.Now(),
	}

	err = userRepoTx.Save(user)
	if err != nil {
		tx.Rollback()
		return nil, common.SystemError
	}

	if req.RoleID > 0 {
		err = userRepoTx.AddUserRole(user.UserID, req.RoleID)
		if err != nil {
			tx.Rollback()
			return nil, common.SystemError
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, common.SystemError
	}

	userResponse := modelToUserResponse(user)
	return &userResponse, nil
}
