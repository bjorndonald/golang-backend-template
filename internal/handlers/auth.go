package handlers

import (
	"fmt"

	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bjorndonald/golang-backend-template/constants"
	"github.com/bjorndonald/golang-backend-template/internal/bootstrap"
	"github.com/bjorndonald/golang-backend-template/internal/helpers"
	"github.com/bjorndonald/golang-backend-template/internal/models"
	"github.com/bjorndonald/golang-backend-template/internal/otp"
	"github.com/gofrs/uuid"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	deps *bootstrap.AppDependencies
}

func NewAuthHandler(
	deps *bootstrap.AppDependencies,
) *AuthHandler {
	return &AuthHandler{
		deps: deps,
	}
}

var (
	constant = constants.New()
)

type ErrorResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

type SuccessResponse struct {
	Message int  `json:"message"`
	Success bool `json:"success"`
}

type RegisterResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    RegisterResponseData `json:"data"`
}

type RegisterResponseData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AuthenticateUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type InputCreateUser struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
}

// LoginResponse represents the response data structure for the login API.
type LoginResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    LoginResponseData `json:"data"`
}

// LoginResponseData represents the data section of the login response.
type LoginResponseData struct {
	JWT string `json:"jwt"`
}

type UpdateAccountInformation struct {
	Country        string `json:"country" validate:"required"`
	Manager        string `json:"manager" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required,numeric"`
	CompanyWebsite string `json:"company_website" validate:"required,url"`
}

// ? EmailInput struct
type EmailInput struct {
	Email string `json:"email" validate:"required"`
}

// ? ResetPasswordInput struct
type ResetPasswordInput struct {
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
}

type OtpVerifyInput struct {
	OTP   string `json:"otp" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type OtpInput struct {
	OTP string `json:"otp" validate:"required"`
}

func checkAgent(userAgent, newAgent models.UserAgent) bool {
	if userAgent.OS == newAgent.OS && userAgent.Platform == newAgent.Platform && userAgent.BrowserName == newAgent.BrowserName && userAgent.Model == newAgent.Model && userAgent.Mobile == newAgent.Mobile {
		return true
	}
	return false
}

func checkLocation(userLocation, newLocation models.GeoLocation) bool {
	if userLocation.City == newLocation.City && userLocation.Country == newLocation.Country && userLocation.Region == newLocation.Region {
		return true
	}
	return false
}

// Authenticate authenticates a user and generates a JWT token.
//
// @Summary Authenticate User
// @Description Authenticate a user by validating their email and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body AuthenticateUser true "User credentials (email and password)"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (a *AuthHandler) Authenticate(c *gin.Context) {
	var input AuthenticateUser
	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(AuthenticateUser)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	user, userExist, err := a.deps.UserRepo.FindByCondition("email = ?", strings.ToLower(input.Email))
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !userExist {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("invalid account credentials"), http.StatusBadRequest)
		return
	}

	clientUrl := constant.ClientUrl

	timeNow, err := helpers.TimeNow("Africa/Lagos")
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !user.EmailVerified {
		helpers.ReturnError(c, "Account not active", fmt.Errorf("account not verified"), http.StatusBadRequest)
		return
	}

	userAgent, _, err := a.deps.AgentRepo.FindByCondition("user_id = ?", user.ID.String())
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("could not get user device info"), http.StatusInternalServerError)
		return
	}

	userLocation, _, err := a.deps.LocationRepo.FindByCondition("user_id = ?", user.ID.String())
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("could not get user location info"), http.StatusInternalServerError)
		return
	}

	loc, agent, err := helpers.GetDeviceLocation(c)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	userLocationCheck := checkLocation(*userLocation, loc)
	if !userLocationCheck {
		baseURL := helpers.GetBaseURL(c)

		user.LastLogin = timeNow

		_, err = a.deps.UserRepo.Save(user)
		if err != nil {
			helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
			return
		}

		go a.deps.EmailService.SendNewLocationEmail(user.FirstName, user.Email, baseURL, loc)
		// helpers.ReturnError(c, "New device needs authorization", fmt.Errorf("New device needs authorization"), http.StatusBadRequest)
		// c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/location", clientUrl))
		// return
	}

	userAgentCheck := checkAgent(*userAgent, agent)
	if !userAgentCheck {
		baseURL := helpers.GetBaseURL(c)

		user.LastLogin = timeNow

		_, err = a.deps.UserRepo.Save(user)
		if err != nil {
			helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
			return
		}

		go a.deps.EmailService.SendNewDeviceEmail(user.FirstName, user.Email, baseURL, agent)
		// helpers.ReturnError(c, "New device needs authorization", fmt.Errorf("New device needs authorization"), http.StatusBadRequest)
		// c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/device", clientUrl))
		// return
	}

	if user.AuthVersion == models.Outdated {
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusFound, fmt.Sprintf("%s/auth/forgot-password", clientUrl))
		return
	}

	hashedPassword := []byte(user.Password)
	plainPassword := []byte(input.Password)
	err = bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)

	if err != nil {
		helpers.ReturnError(c, "Email and Password is not correct", err, http.StatusBadRequest)
		return
	}

	user.LastLogin = timeNow
	user.IP = c.ClientIP()

	_, err = a.deps.UserRepo.Save(user)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	accessToken, err := helpers.GenerateAccessToken(constant.JWTSecretKey, user.Email, user.FirstName, (user.ID).String())
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	go a.deps.EmailService.SendOTPEmail(user.FirstName, user.Email)

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusFound, fmt.Sprintf("%s/auth/2fa?token=%s", clientUrl, accessToken))
}

// LogOut is a route handler that logs out the user.
//
// This endpoint is used to log out the user.
//
// @Summary Log out User
// @Description Log out User
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {string} string "Returns 'success' "
// @Failure 400 {string} string "Returns error message"
// @Router /auth/2fa/{token} [post]
func (a *AuthHandler) LogOut(c *gin.Context) {

	c.SetCookie("refreshToken", "", 0, "/", "", true, true)
	helpers.ReturnJSON(c, "OTP sent successfully", nil, http.StatusOK)
}

// Send2faEmail is a route handler that send an otp to user's email address.
//
// This endpoint is used to send an otp code to the user's email address.
//
// @Summary Send 2fa OTP
// @Description Send 2fa OTP
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body Email true "Email"
// @Success 200 {string} string "Returns 'success' "
// @Failure 400 {string} string "Returns error message"
// @Router /auth/2fa/{token} [post]
func (a *AuthHandler) Send2FAEmail(c *gin.Context) {
	accessToken := c.Params.ByName("token")

	token, err := jwt.ParseWithClaims(
		accessToken, &helpers.AuthTokenJwtClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(constant.JWTSecretKey), nil
		})

	claims := token.Claims.(*helpers.AuthTokenJwtClaim)

	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	user, _, err := a.deps.UserRepo.FindByCondition("id", claims.UserId)

	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	go a.deps.EmailService.SendOTPEmail(user.FirstName, user.Email)

	helpers.ReturnJSON(c, "OTP sent successfully", nil, http.StatusOK)
}

// Verify2faEmail is a route handler that verifies the user's email address.
//
// This endpoint is used to verify the user's email address by providing the email and OTP token.
//
// @Summary Verify email address
// @Description Verifies the user's email address
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body OtpVerifyInput true "Email and OTP"
// @Success 200 {string} string "Returns 'success' "
// @Failure 400 {string} string "Returns error message"
// @Router /auth/2fa/verify/{token} [post]
func (a *AuthHandler) Verify2FAEmail(c *gin.Context) {
	token := c.Params.ByName("token")
	var input OtpInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Error parsing request", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(OtpInput)
	if !ok {
		helpers.ReturnError(c, "Error parsing request", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	authToken, err := jwt.ParseWithClaims(
		token, &helpers.AuthTokenJwtClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(constant.JWTSecretKey), nil
		})

	claims := authToken.Claims.(*helpers.AuthTokenJwtClaim)

	if err != nil {
		helpers.ReturnError(c, "Access token is not valid", err, http.StatusInternalServerError)
		return
	}

	valid := otp.OTPManage.VerifyOTP(claims.Email, input.OTP)

	if !valid {
		helpers.ReturnJSON(c, "OTP not valid", nil, http.StatusBadRequest)
		return
	}

	accessToken, err := helpers.GenerateAccessToken(constant.JWTSecretKey, claims.Email, claims.Name, claims.UserId)
	if err != nil {
		helpers.ReturnError(c, "Could not generate token", err, http.StatusInternalServerError)
		return
	}

	refreshToken, err := helpers.GenerateRefreshToken(constant.JWTSecretKey, claims.Email, claims.Name, claims.UserId)
	if err != nil {
		helpers.ReturnError(c, "Could not generate token", err, http.StatusInternalServerError)
		return
	}

	c.SetCookie("refreshToken", refreshToken, 60*60, "/", "", true, true)

	helpers.ReturnJSON(c, "Authenticated successfully", map[string]interface{}{
		"access_token": accessToken,
		"expires_in":   time.Now().Local().Add(time.Minute * 15),
	}, http.StatusOK)
}

// Authenticate authenticates a user and generates a JWT token.
//
// @Summary Authenticate User
// @Description Authenticate a user by validating their email and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body AuthenticateUser true "User credentials (email and password)"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh-token [post]
func (a *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil {
		helpers.ReturnError(c, "Refresh token missing", err, http.StatusInternalServerError)
		return
	}

	claims, err := helpers.ValidateToken(refreshToken)
	if err != nil {
		helpers.ReturnError(c, "Invalid refresh token", err, http.StatusInternalServerError)
		return
	}

	accessToken, err := helpers.GenerateAccessToken(constant.JWTSecretKey, claims.Email, claims.Name, claims.UserId)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	refreshToken, err = helpers.GenerateRefreshToken(constant.JWTSecretKey, claims.Email, claims.Name, claims.UserId)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	c.SetCookie("refreshToken", refreshToken, 60*60, "/", "", true, true)

	helpers.ReturnJSON(c, "Refreshed successfully", map[string]interface{}{
		"access_token": accessToken,
		"expires_in":   time.Now().Local().Add(time.Minute * 15),
	}, http.StatusOK)
}

func (a *AuthHandler) findUserOrError(email string) (user *models.User, err error) {
	user, userExist, err := a.deps.UserRepo.FindByCondition("email = ?", email)
	if err != nil {
		return nil, err
	}
	if !userExist {
		return nil, helpers.NewError("user not found")
	}
	return user, nil
}

// Register creates a new user account.
//
// @Summary Create a new user
// @Description Create a new user account with the provided information
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body InputCreateUser true "User data to create an account"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (a *AuthHandler) CreateUser(c *gin.Context) {
	var input InputCreateUser
	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Invalid request", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(InputCreateUser)
	if !ok {
		helpers.ReturnError(c, "Invalid request", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	_, found, err := a.deps.UserRepo.FindByCondition("email", input.Email)
	if err != nil {
		helpers.ReturnError(c, "error getting user", err, http.StatusInternalServerError)
		return
	}

	if found {
		helpers.ReturnError(c, "User already found", fmt.Errorf("account already exists"), http.StatusConflict)
		return
	}

	if input.Password != input.ConfirmPassword {
		helpers.ReturnError(c, "Passwords don't match", fmt.Errorf("passwords don't match"), http.StatusConflict)
		return
	}

	hash, err := helpers.HashPassword(input.Password)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	userID, err := uuid.NewV7()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	// create record
	user := &models.User{
		Email:         strings.ToLower(input.Email),
		Password:      hash,
		ID:            userID,
		IP:            c.ClientIP(),
		Role:          models.UserRole,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		Status:        models.InactiveAccount,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
	}

	if err := a.deps.UserRepo.Create(user); err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	baseURL := helpers.GetBaseURL(c)

	go a.deps.EmailService.SendNewUserEmail(user.FirstName, user.Email, baseURL)

	eventJSON, err := json.Marshal(user)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}
	go a.deps.EventProducer.BroadCast(1, "signup", eventJSON)

	helpers.ReturnJSON(c, "Account created successfully", user, http.StatusCreated)
}

// VerifyEmail is a route handler that verifies the user's email address.
//
// This endpoint is used to verify the user's email address by providing the email and OTP token.
//
// @Summary Verify email address
// @Description Verifies the user's email address
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email path string true "User's email address"
// @Param otp path string true "One-time password (OTP) token"
// @Success 302 {string} string "Redirects to the client URL with a jwt token"
// @Failure 302 {string} string "Redirects to the client URL with an error code"
// @Router /auth/verify/{email}/{otp} [get]
func (a *AuthHandler) VerifyEmail(c *gin.Context) {
	email := c.Param("email")
	token := c.Param("otp")

	user, userExist, err := a.deps.UserRepo.FindByCondition("email = ?", email)
	clientUrl := constant.ClientUrl

	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/signin?error=500", clientUrl))
		return
	}
	if !userExist {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/signin?error=402", clientUrl))
		return
	}

	location, agent, err := helpers.GetDeviceLocation(c)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	agentID, err := uuid.NewV7()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	agent.ID = agentID
	agent.UserID = user.ID

	if err := a.deps.AgentRepo.Create(&agent); err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	locationID, err := uuid.NewV7()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	userLocation := &models.GeoLocation{
		UserID:   user.ID,
		ID:       locationID,
		IP:       location.IP,
		City:     location.City,
		Country:  location.Country,
		Region:   location.Region,
		Location: location.Location,
	}

	if err := a.deps.LocationRepo.Create(userLocation); err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	accessToken, err := helpers.GenerateAccessToken(constant.JWTSecretKey, user.Email, user.FirstName, user.ID.String())

	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/signin?error=500", clientUrl))
		return
	}

	refreshToken, err := helpers.GenerateRefreshToken(constant.JWTSecretKey, user.Email, user.FirstName, user.ID.String())

	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/signin?error=500", clientUrl))
		return
	}

	if user.EmailVerified {
		c.SetCookie("refreshToken", refreshToken, 60*60, "/", "", true, true)
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/signin?&access_token=%s", clientUrl, accessToken))
		return
	}

	valid := otp.OTPManage.VerifyOTP(email, token)

	if !valid {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/signin?error=401V", clientUrl))
		return
	}

	user.EmailVerified = true
	user.Status = models.ActiveAccount
	user.UpdatedAt = time.Now()

	_, err = a.deps.UserRepo.Save(user)

	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/signin?error=500", clientUrl))
		return
	}

	c.SetCookie("refreshToken", refreshToken, 60*60, "/", "", true, true)
	c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/signin?&access_token=%s", clientUrl, accessToken))
}

// ForgotPassword is a route handler that sends the reset otp to the user's email address.
//
// This endpoint is used to send the otp to the user's email address by providing the email.
//
// @Summary Sends reset OTP
// @Description Sends the reset OTP to the user's email address
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body EmailInput true "Input (email)"
// @Success 200 {string} string "Returns 'success' "
// @Failure 400 {string} string "Returns error message"
// @Router /auth/forgot-password [post]
func (a *AuthHandler) ForgotPassword(c *gin.Context) {
	var input EmailInput
	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(EmailInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	userFound, err := a.findUserOrError(input.Email)

	if userFound == nil && err != nil {
		helpers.ReturnJSON(c, "User not found", err, http.StatusBadRequest)
		return
	}

	var fullName string = userFound.FirstName + userFound.LastName
	var email string = userFound.Email

	if strings.Contains(fullName, " ") {
		fullName = strings.Split(fullName, " ")[1]
	}
	a.deps.EmailService.SendForgotPasswordEmail(fullName, email)
	clientUrl := constant.ClientUrl

	// helpers.ReturnJSON(c, "Action successful", nil, http.StatusOK)
	c.JSON(http.StatusFound, fmt.Sprintf("%s/auth/reset-password?email=%s", clientUrl, email))
}

// VerifyResetOTP is a route handler that verifies the user's email address.
//
// This endpoint is used to verify the user's email address by providing the email and OTP token.
//
// @Summary Verify email address
// @Description Verifies the user's email address
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body OtpVerifyInput true "Input (token and email)"
// @Success 200 {string} string "Returns 'success and JWT' "
// @Failure 400 {string} string "Returns error message"
// @Router /auth/forgot-password/verify [post]
func (a *AuthHandler) VerifyResetOTP(c *gin.Context) {
	var input OtpVerifyInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Invalid input", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(OtpVerifyInput)
	if !ok {
		helpers.ReturnError(c, "Error validating input", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	user, userExist, err := a.deps.UserRepo.FindByCondition("email", input.Email)

	if err != nil {
		helpers.ReturnError(c, "Could not get user", err, http.StatusInternalServerError)
		return
	}

	if !userExist {
		helpers.ReturnError(c, "User does not exist", fmt.Errorf("user account does not exist"), http.StatusBadRequest)
		return
	}

	valid := otp.OTPManage.VerifyOTP(input.Email, input.OTP)

	if !valid {
		helpers.ReturnError(c, "OTP not valid", fmt.Errorf("invalid opt"), http.StatusBadRequest)
		return
	}

	jwtToken, err := helpers.GenerateAccessToken(constant.JWTSecretKey, user.Email, user.FirstName, user.ID.String())

	if err != nil {
		helpers.ReturnError(c, "Error generating access token", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusFound, fmt.Sprintf("%s/auth/change-password?reset_token=%s", constant.ClientUrl, jwtToken))
}

// ResetPassword is a route handler for resetting the user's password.
// It requires a valid JWT token and a JSON request body with new credentials.
//
// @Summary Reset the user's password
// @Description Reset the user's password using a JWT token and new credentials.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param reset-token path string true "JWT token for resetting the password"
// @Param credentials body ResetPasswordInput true "New password and password confirmation"
// @Success 200 {string} string "Success: Password reset"
// @Failure 400 {string} string "Error: Invalid input or token"
// @Router /auth/reset-password/confirm/{reset-token} [post]
func (a *AuthHandler) ResetPassword(c *gin.Context) {
	resetToken := c.Params.ByName("reset-token")

	var input ResetPasswordInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(ResetPasswordInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	token, err := jwt.ParseWithClaims(
		resetToken, &helpers.AuthTokenJwtClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(constant.JWTSecretKey), nil
		})

	claims := token.Claims.(*helpers.AuthTokenJwtClaim)

	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	user, _, err := a.deps.UserRepo.FindByCondition("id", claims.UserId)

	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if user == nil {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("user not found"), http.StatusNotFound)
		return
	}

	if input.Password != input.PasswordConfirm {
		helpers.ReturnError(c, "Passwords do not match", fmt.Errorf("passwords do not match"), http.StatusBadRequest)
		return
	}

	oldPassword := []byte(user.Password)
	plainPassword := []byte(input.Password)
	err = bcrypt.CompareHashAndPassword(oldPassword, plainPassword)
	if err == nil {
		helpers.ReturnError(c, "Please input a different password than one used before.", fmt.Errorf("Password is the same"), http.StatusBadRequest)
		return
	}

	hashedPassword, _ := helpers.HashPassword(input.Password)
	user.Password = hashedPassword
	user.EmailVerified = true
	user.AuthVersion = models.UpToDate
	user.Status = models.ActiveAccount
	user.UpdatedAt = time.Now()

	_, err = a.deps.UserRepo.Save(user)

	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Password updated successfully", user, http.StatusOK)
}
