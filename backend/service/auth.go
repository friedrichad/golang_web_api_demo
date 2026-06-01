package service

import (
	"encoding/base64"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/friedrichad/golang_web_api_demo/backend/common"
	"github.com/friedrichad/golang_web_api_demo/backend/middleware"
	"github.com/friedrichad/golang_web_api_demo/backend/model"
	"github.com/friedrichad/golang_web_api_demo/backend/rabbitmq"
	"github.com/friedrichad/golang_web_api_demo/backend/redis"
	"github.com/friedrichad/golang_web_api_demo/backend/repository"
	"github.com/friedrichad/golang_web_api_demo/backend/shared"
	"github.com/friedrichad/golang_web_api_demo/backend/utils"
	"github.com/pquerna/otp/totp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Authentication(c *gin.Context) (*model.TokenResponse, *common.Error)
	Register(c *gin.Context) (*model.UserResponse, *common.Error)
	Logout(c *gin.Context) *common.Error
	Setup2FA(c *gin.Context) (*model.TwoFASetupResponse, *common.Error)
	VerifySetup2FA(c *gin.Context) (*model.TwoFAVerifyResponse, *common.Error)
	Verify2FA(c *gin.Context) (*model.TokenResponse, *common.Error)
}

type AuthService struct {
	rabbitMQ            *rabbitmq.RabbitMQ
	repository          repository.IUserRepository
	accessTokenExpired  int
	refreshTokenExpired int
	clientId            string
	basicAuth           string
	jwtSecret           string
	clientTypeNoExp     []string
}

func NewAuthService() IAuthService {
	rmq, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		log.Println(
			"RabbitMQ unavailable:",
			err,
		)
	}
	return &AuthService{
		rabbitMQ:            rmq,
		repository:          repository.NewUserRepository(),
		accessTokenExpired:  viper.GetInt("oauth.access-token-expired"),
		refreshTokenExpired: viper.GetInt("oauth.refresh-token-expired"),
		clientId:            viper.GetString("oauth.client-id"),
		basicAuth:           base64.StdEncoding.EncodeToString([]byte(viper.GetString("oauth.client-id") + ":" + viper.GetString("oauth.client-secret"))),

		jwtSecret: viper.GetString(
			"oauth.jwt-secret",
		),

		clientTypeNoExp: viper.GetStringSlice(
			"oauth.client_type_no_expire",
		),
	}
}

func (a *AuthService) Authentication(c *gin.Context) (*model.TokenResponse, *common.Error) {
	grantType := c.Request.FormValue("grant_type")
	if grantType == "password" {
		user, err := a.ValidateUser(c)
		if err != nil {
			return nil, err
		}
		if user.TwoFaEnabled {
			// Return user info without access token when 2FA is required
			response := &model.TokenResponse{
				Required2FA: true,
				UserId:      user.UserID,
			}
			response.Id = strconv.FormatInt(int64(user.UserID), 10)
			response.Username = user.Username
			response.Active = true

			// Fetch Position
			positionRepo := repository.NewPositionRepository()
			position, err := positionRepo.GetPositionById(user.PositionID)
			if err == nil && position != nil {
				response.PositionName = position.PositionName
				response.Level = position.PositionLevel
			} else {
				response.PositionName = ""
				response.Level = 999
			}

			authorities, err := a.repository.GetAuthorities(user.UserID)
			if err != nil && err.Error() != "record not found" {
				return nil, common.SystemError
			}
			response.Authorities = authorities

			return response, nil
		}

		return createNewToken(user, a, c)
	}
	if grantType == "refresh_token" {
		return refreshToken(c, a)
	}
	return nil, &common.Error{Code: "400", Message: "Grant type không hỗ trợ"}
}

func (a *AuthService) ValidateUser(c *gin.Context) (*model.User, *common.Error) {
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

	if user == nil {
		return nil, common.AuthenticationFail
	}

	if IsLockedAccount(user.UserID, 5) {
		log.Printf("Tài khoản user_id %d đang bị khóa do đăng nhập sai quá nhiều lần", user.UserID)
		return nil, common.AccountLocked
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		locked, count := LockAccount(user.UserID, 5)
		if locked {
			log.Printf("Tài khoản user_id %d bị khóa do đăng nhập sai quá nhiều lần (%d lần)", user.UserID, count)
			return nil, common.AccountLocked
		}
		return nil, common.AuthenticationFail
	}

	return user, nil
}

func createNewToken(user *model.User, a *AuthService, c *gin.Context) (*model.TokenResponse, *common.Error) {
	sessionId := utils.GetOrCreateSessionID(c, time.Until(time.Unix(getExpiredTime(a.accessTokenExpired), 0)))
	if existedUserId, err := utils.BrowserHasSession(sessionId); err == nil && existedUserId != "" {
		log.Printf("Trình duyệt đã có session_id %s với user_id %s", sessionId, existedUserId)
		return nil, common.AlreadyLoggedIn
	}

	response := &model.TokenResponse{
		AccessToken: "",
		TokenType:   "bearer",
	}
	log.Printf("UserId: %v", user.UserID)
	response.Id = strconv.FormatInt(int64(user.UserID), 10)
	response.Username = user.Username
	response.Active = true
	// Fetch Position
	positionRepo := repository.NewPositionRepository()
	position, err := positionRepo.GetPositionById(user.PositionID)
	if err == nil && position != nil {
		response.PositionName = position.PositionName
		response.Level = position.PositionLevel
	} else {
		response.PositionName = ""
		response.Level = 999
	}
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
	err = utils.SaveBrowserSession(sessionId, response.Id, ttl)
	if err != nil {
		return nil, common.SystemError
	}
	err = redis.Save(redis.Rdb, "auth:token:"+response.Id, response.AccessToken, ttl)
	if err != nil {
		return nil, common.SystemError
	}
	userInfo := shared.UserInfo{
		UserId:   user.UserID,
		UserName: user.Username,
		PositionInfo: shared.PositionInfo{
			PositionId:    user.PositionID,
			PositionName:  response.PositionName,
			PositionLevel: response.Level,
		},
		IsOP: response.IsOP,
	}
	err = redis.SaveUserInfoCache(redis.Rdb, userInfo, ttl)
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
	claims["position_id"] = token.PositionID
	claims["position_name"] = token.PositionName
	claims["position_level"] = token.Level
	claims["is_op"] = token.IsOP
	accessToken, err := t.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func refreshToken(c *gin.Context, a *AuthService) (*model.TokenResponse, *common.Error) {
	sessionID, err := c.Cookie("session_id")
	if err != nil || sessionID == "" {
		return nil, common.TokenInvalid
	}
	userID, err := utils.BrowserHasSession(sessionID)
	if err != nil || userID == "" {
		return nil, common.TokenInvalid
	}
	claims, ok := extractClaims(c.Request.FormValue("refresh_token"), []byte(a.jwtSecret))
	if !ok {
		return nil, common.TokenInvalid
	}
	if claims.Id != userID {
		log.Printf("Session userID %s không khớp với token userID %s", userID, claims.Id)
		return nil, common.TokenInvalid
	}
	response := &model.TokenResponse{
		AccessToken: "",
		TokenType:   "bearer",
	}
	response.Id = claims.Id
	response.Username = claims.Username
	response.PositionID = claims.PositionID
	response.PositionName = claims.PositionName
	response.Level = claims.Level
	response.IsOP = claims.IsOP
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
	// Check if token matches the latest token stored in Redis (not just existence)
	storedToken, err := redis.Get(redis.Rdb, "auth:token:"+claims.Id)
	if err != nil {
		log.Printf("Error getting token from Redis: %v", err)
		return nil, false
	}
	if storedToken == "" || storedToken != tokenStr {
		log.Printf("Token does not match the latest token in Redis or token has been revoked")
		return nil, false
	}
	return claims, true
}

func (a *AuthService) Register(c *gin.Context) (*model.UserResponse, *common.Error) {
	var req model.RegisterRequest
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

	userRepo := a.repository.(*repository.UserRepository)

	user := &model.User{
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		StatusInt:    1,
		CreatedAt:    time.Now(),
		PositionID:   req.PositionID,
	}

	err = userRepo.Save(user)
	if err != nil {
		return nil, common.SystemError
	}

	if a.rabbitMQ != nil {

		err := a.rabbitMQ.Publish(
			"user.created",
			map[string]interface{}{
				"user_id": user.UserID,
				"email":   user.Email,
			},
		)

		if err != nil {
			log.Println(
				"Publish user.created failed:",
				err,
			)
		}
	}
	userResponse := modelToUserResponse(user)
	return &userResponse, nil
}

func (a *AuthService) Logout(c *gin.Context) *common.Error {
	token := c.GetHeader("Authorization")
	if token == "" {
		return common.TokenInvalid
	}
	token = strings.Replace(token, "Bearer ", "", 1)
	claims := &model.Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})

	if err != nil && err.Error() != "Token is expired" {
		log.Printf("Failed to parse token: %v", err)
		return common.TokenInvalid
	}
	ttl := time.Until(time.Unix(claims.RefreshExp, 0))
	if ttl <= 0 {
		ttl = 1 * time.Second // Minimum TTL
	}
	err = redis.AddToBlacklist(token, ttl)
	if err != nil {
		log.Printf("Failed to add token to blacklist: %v", err)
		return common.SystemError
	}
	err = redis.Delete(redis.Rdb, "auth:token:"+claims.Id)
	if err != nil {
		log.Printf("Failed to delete token from Redis: %v", err)
		return common.SystemError
	}
	sessionID, _ := c.Cookie("session_id")
	if sessionID != "" {
		redis.Delete(redis.Rdb, "auth:browser:"+sessionID)
	}
	utils.ClearSessionCookie(c)
	return nil
}

func LockAccount(userID int, attempt int) (bool, int) {
	key := "auth:lock:" + strconv.Itoa(userID)
	value, _ := redis.Get(redis.Rdb, key)
	count := 0
	if value != "" {
		count, _ = strconv.Atoi(value)
	}
	count++
	err := redis.Save(redis.Rdb, key, strconv.Itoa(count), 15*time.Minute)
	if err != nil {
		log.Printf("Redis save error: %v", err)
	}
	if count >= attempt {
		return true, count
	}
	return false, count
}

func IsLockedAccount(userID int, attempt int) bool {
	key := "auth:lock:" + strconv.Itoa(userID)
	value, _ := redis.Get(redis.Rdb, key)
	if value == "" {
		return false
	}
	count, _ := strconv.Atoi(value)
	if count >= attempt {
		return true
	}
	return false
}

func (a *AuthService) Setup2FA(c *gin.Context) (*model.TwoFASetupResponse, *common.Error) {

	userId, err := strconv.Atoi(middleware.GetUserID(c))
	if err != nil {
		log.Printf("Invalid user id: %s", err.Error())
		return nil, common.RequestInvalid
	}

	user, err := a.repository.GetById(userId)
	if err != nil {
		log.Printf("Error when get user by id: %s", err.Error())
		return nil, common.NotFound
	}
	key, err := Generate2FA(user.Email)
	if err != nil {
		log.Printf("Error when generate 2FA key: %s", err.Error())
		return nil, common.SystemError
	}
	user.TwoFaSecret = key.Secret()
	err = a.repository.Update(user)
	if err != nil {
		log.Printf("Error when update user with 2Fa secret: %s", err.Error())
		return nil, common.SystemError
	}
	return &model.TwoFASetupResponse{
		Secret: key.Secret(),
		QrUrl:  key.URL(),
	}, nil
}

func (a *AuthService) VerifySetup2FA(c *gin.Context) (*model.TwoFAVerifyResponse, *common.Error) {

	userId, err := strconv.Atoi(middleware.GetUserID(c))
	if err != nil {
		return nil, common.NotFound
	}

	user, err := a.repository.GetById(userId)
	if err != nil {
		return nil, common.NotFound
	}

	var body model.TwoFAVerifyRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, common.RequestInvalid
	}

	valid := totp.Validate(body.Code, user.TwoFaSecret)

	if !valid {
		return nil, &common.Error{
			Code:    "400",
			Message: "Invalid 2FA code",
		}
	}

	user.TwoFaEnabled = true

	err = a.repository.Update(user)
	if err != nil {
		return nil, common.SystemError
	}

	return &model.TwoFAVerifyResponse{
		Message: "2FA enabled",
	}, nil
}

func (a *AuthService) Verify2FA(c *gin.Context) (*model.TokenResponse, *common.Error) {

	var body struct {
		UserId int    `json:"user_id"`
		Code   string `json:"code"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, common.RequestInvalid
	}

	user, err := a.repository.GetById(body.UserId)
	if err != nil {
		return nil, common.NotFound
	}
	if !user.TwoFaEnabled {
		return nil, &common.Error{
			Code:    "400",
			Message: "2FA chưa được bật",
		}
	}

	valid := totp.Validate(body.Code, user.TwoFaSecret)

	if !valid {
		return nil, &common.Error{
			Code:    "401",
			Message: "Mã xác thực không đúng",
		}
	}
	return createNewToken(user, a, c)
}
