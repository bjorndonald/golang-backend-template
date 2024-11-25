package routes

import (
	"github.com/bjorndonald/golang-backend-template/internal/bootstrap"
	"github.com/bjorndonald/golang-backend-template/internal/handlers"
	"github.com/bjorndonald/golang-backend-template/internal/middleware"
	"github.com/bjorndonald/golang-backend-template/internal/validators"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, d *bootstrap.AppDependencies) {
	userRouter := router.Group("user")

	handler := handlers.NewUserHandler(d)

	userRouter.GET("/profile", middleware.JWTMiddleware(d.DatabaseService), handler.UserProfile)
	userRouter.PUT("/profile", middleware.JWTMiddleware(d.DatabaseService), validators.ValidateUpdateUserProfile, handler.UpdateUserProfile)
	userRouter.PUT("/photo", middleware.JWTMiddleware(d.DatabaseService), middleware.CloudinaryUploadMiddleware(), handler.UpdateUserPhoto)

	// OTP

	userRouter.POST("/otp", middleware.JWTMiddleware(d.DatabaseService), handler.SendOTP)
	userRouter.POST("/otp/verify", middleware.JWTMiddleware(d.DatabaseService),
		validators.ValidateOTPSchema, handler.VerifyOTP)
}
