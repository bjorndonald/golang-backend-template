package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/bjorndonald/golang-backend-template/constants"
	"github.com/bjorndonald/golang-backend-template/internal/bootstrap"
	"github.com/bjorndonald/golang-backend-template/internal/helpers"
	"github.com/bjorndonald/golang-backend-template/internal/manager"
	"github.com/bjorndonald/golang-backend-template/internal/models"
	"github.com/bjorndonald/golang-backend-template/internal/otp"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	deps          *bootstrap.AppDependencies
	streamManager *manager.Manager
}

func NewUserHandler(deps *bootstrap.AppDependencies,
) *UserHandler {
	return &UserHandler{
		deps:          deps,
		streamManager: manager.NewGameManager(),
	}
}

type FileUploadResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type UserProfile struct {
	Email     string    `json:"email"`
	Userid    string    `json:"userid"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateUserProfileInput struct {
	Email       string `json:"email" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Bio         string `json:"bio" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type UpdateRoleInput struct {
	Role string `json:"role" validate:"required"`
}

type UserListResponse struct {
	Message int               `json:"message"`
	Success bool              `json:"success"`
	Data    []models.UserInfo `json:"data"`
}

type UserResponse struct {
	Message int             `json:"message"`
	Success bool            `json:"success"`
	Data    models.UserInfo `json:"data"`
}

// UpdateUserProfile is a route handler that handles updating the user profile
//
// # This endpoint is used to update the user profile
//
// @Summary Update user profile
// @Description Updates some details about the user
// @Tags User
// @Accept json
// @Produce json
// @Param credentials body UpdateUserProfileInput true "update user profile"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /user [put]
func (u *UserHandler) UpdateUserProfile(c *gin.Context) {
	var input UpdateUserProfileInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(UpdateUserProfileInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	claims, err := helpers.GetAuthenticatedUser(c)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	user, found, err := u.deps.UserRepo.FindByCondition("email = ?", claims.Email)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("user not found"), http.StatusNotFound)
		return
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Email = input.Email
	user.Bio = input.Bio
	user.PhoneNumber = input.PhoneNumber

	_, err = u.deps.UserRepo.Save(user)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Profile updated successfully", user, http.StatusOK)
}

// UpdateUserPhoto is a route handler that handles updating the user photo
//
// # This endpoint is used to update the user photo
//
// @Summary Update user photo
// @Description Updates some details about the user
// @Tags User
// @Accept json
// @Produce json

// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /user/photo [put]
func (u *UserHandler) UpdateUserPhoto(c *gin.Context) {
	image, exists := c.Get("image")

	if !exists {
		helpers.ReturnError(c, "Couldn't get image", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	imageUrl, ok := image.(string)
	if !ok {
		helpers.ReturnError(c, "Couldn't get image", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	claims, err := helpers.GetAuthenticatedUser(c)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	user, found, err := u.deps.UserRepo.FindByCondition("email = ?", claims.Email)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("user not found"), http.StatusNotFound)
		return
	}

	user.Photo = imageUrl

	_, err = u.deps.UserRepo.Save(user)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Profile updated successfully", user, http.StatusOK)
}

// UserProfile is a route handler that retrieves the user profile of the authenticated user.
//
// This endpoint is used to get the profile information of the authenticated user based on the JWT claims.
//
// @Summary Get user profile
// @Description Retrieves the profile information of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse
// @Router /user/profile [get]
func (u *UserHandler) UserProfile(c *gin.Context) {

	claimsRaw, exists := c.Get("claims")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	authClaims, ok := claimsRaw.(*helpers.AuthTokenJwtClaim)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, f, err := u.deps.UserRepo.FindByCondition("id = ?", authClaims.UserId)

	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !f {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("account not found"), http.StatusNotFound)
		return
	}

	helpers.ReturnJSON(c, "Profile retrieved", user, http.StatusOK)

}

// SendOTP is a route handler that send an otp to user's email address.
//
// This endpoint is used to send an otp code to the user's email address.
//
// @Summary Send 2fa OTP
// @Description Send 2fa OTP
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/otp [post]
func (a *UserHandler) SendOTP(c *gin.Context) {
	claims, err := helpers.GetAuthenticatedUser(c)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	user, found, err := a.deps.UserRepo.FindByCondition("email = ?", claims.Email)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("user not found"), http.StatusNotFound)
		return
	}

	go a.deps.EmailService.SendOTPEmail(user.FirstName, user.Email)

	helpers.ReturnJSON(c, "OTP sent successfully", nil, http.StatusOK)
}

// VerifyOTP is a route handler that verifies the user's email address.
//
// This endpoint is used to verify the user's email address by providing the email and OTP token.
//
// @Summary Verify email address
// @Description Verifies the user's email address
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body OtpInput true "OTP"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/otp/verify [post]
func (a *UserHandler) VerifyOTP(c *gin.Context) {
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

	claims, err := helpers.GetAuthenticatedUser(c)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	user, found, err := a.deps.UserRepo.FindByCondition("email = ?", claims.Email)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("user not found"), http.StatusNotFound)
		return
	}

	valid := otp.OTPManage.VerifyOTP(user.Email, input.OTP)

	if !valid {
		helpers.ReturnJSON(c, "OTP not valid", nil, http.StatusBadRequest)
		return
	}

	helpers.ReturnJSON(c, "OTP is valid", nil, http.StatusOK)
}

func uploadFile(file *multipart.FileHeader, folderName string) (resp *uploader.UploadResult, err error) {

	env_ := constants.New()
	// Open the uploaded file
	fileOpened, err := file.Open()
	if err != nil {
		// Handle error
		return nil, err
	}

	defer fileOpened.Close()
	url := fmt.Sprintf("cloudinary://%s:%s@%s", env_.CloudinaryAPIKey, env_.CloudinaryApiSecret, env_.CloudinaryName)

	cld, err := cloudinary.NewFromURL(url)

	if err != nil {
		return nil, err
	}
	var ctx = context.Background()
	// Upload the image to Cloudinary
	resp, err = cld.Upload.Upload(ctx, fileOpened, uploader.UploadParams{PublicID: file.Filename,
		Folder: folderName,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// NotFound returns custom 404 page
func NotFound(c *gin.Context) {
	c.Status(404)
	c.File("./static/private/404.html")
}
