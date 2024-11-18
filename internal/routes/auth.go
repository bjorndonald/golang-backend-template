package routes

import (
	"github.com/bjorndonald/lasgcce/internal/bootstrap"
	"github.com/bjorndonald/lasgcce/internal/handlers"
	"github.com/bjorndonald/lasgcce/internal/validators"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, d *bootstrap.AppDependencies) {

	handler := handlers.NewAuthHandler(d)

	authRouter := router.Group("/auth")

	authRouter.POST("/register", validators.ValidateRegisterUserSchema, handler.Register)
	authRouter.POST("/login", validators.ValidateLoginUser, handler.Authenticate)
	authRouter.GET("/verify/:email/:otp", handler.VerifyEmail)

	authRouter.POST("/forgot-password/verify", validators.ValidateResetOTPVerifySchema, handler.VerifyResetOTP)
	authRouter.POST("/forgot-password", validators.ValidateResetUserSchema, handler.ForgotPassword)
	authRouter.POST("/reset-password/confirm/:reset-token", validators.ValidateResetPasswordSchema, handler.ResetPassword)

}
