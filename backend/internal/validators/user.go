package validators

import (
	"net/http"

	"github.com/bjorndonald/golang-backend-template/internal/handlers"
	"github.com/bjorndonald/golang-backend-template/internal/helpers"
	"github.com/bjorndonald/golang-backend-template/validator"
	"github.com/gin-gonic/gin"
)

func ValidateRegisterUserSchema(c *gin.Context) {
	var body handlers.InputCreateUser
	bindAndValidate(c, &body)
	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateAccountResetScheme(c *gin.Context) {
	var body helpers.EmailInput
	bindAndValidate(c, &body)
	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateOTPSchema(c *gin.Context) {
	var body handlers.OtpInput
	bindAndValidate(c, &body)
	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateOTPVerifySchema(c *gin.Context) {
	var body handlers.OtpVerifyInput
	bindAndValidate(c, &body)
	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateLoginUser(c *gin.Context) {
	var body handlers.AuthenticateUser
	bindAndValidate(c, &body)
	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateUpdateUserProfile(c *gin.Context) {
	var body handlers.UpdateUserProfileInput
	bindAndValidate(c, &body)
	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateResetUserSchema(c *gin.Context) {
	var body handlers.EmailInput
	bindAndValidate(c, &body)
	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateResetPasswordSchema(c *gin.Context) {
	var body handlers.ResetPasswordInput
	bindAndValidate(c, &body)
	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateRoleSchema(c *gin.Context) {
	var body handlers.UpdateRoleInput
	bindAndValidate(c, &body)
	c.Set("validatedRequestBody", body)
	c.Next()
}

func bindAndValidate(c *gin.Context, body interface{}) {
	if err := c.ShouldBindJSON(body); err != nil {
		helpers.ReturnError(c, "Error validating input", err, http.StatusBadRequest)
		c.Abort()
		return
	}

	if err := validator.Validate(body); err != nil {
		helpers.ReturnError(c, "Error validating input", err, http.StatusBadRequest)
		c.Abort()
		return
	}
}
