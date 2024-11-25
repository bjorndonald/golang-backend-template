package routes

import (
	"github.com/bjorndonald/golang-backend-template/internal/bootstrap"
	"github.com/bjorndonald/golang-backend-template/internal/handlers"
	"github.com/bjorndonald/golang-backend-template/internal/validators"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, d *bootstrap.AppDependencies) {

	handler := handlers.NewAuthHandler(d)

	authRouter := router.Group("/auth")

	authRouter.POST("/register", validators.ValidateRegisterUserSchema, handler.CreateUser)
	authRouter.POST("/login", validators.ValidateLoginUser, handler.Authenticate)
	authRouter.POST("/logout", handler.LogOut)
	authRouter.POST("/refresh-token", handler.RefreshToken)
	authRouter.GET("/verify/:email/:otp", handler.VerifyEmail)

	authRouter.POST("/forgot-password/verify", validators.ValidateOTPVerifySchema, handler.VerifyResetOTP)
	authRouter.POST("/forgot-password", validators.ValidateResetUserSchema, handler.ForgotPassword)
	authRouter.POST("/reset-password/confirm/:reset-token", validators.ValidateResetPasswordSchema, handler.ResetPassword)

	authRouter.POST("/2fa/:token", handler.Send2FAEmail)
	authRouter.POST("/2fa/verify/:token", validators.ValidateOTPSchema, handler.Verify2FAEmail)

}
