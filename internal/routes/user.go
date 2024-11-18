package routes

import (
	"github.com/bjorndonald/lasgcce/internal/bootstrap"
	"github.com/bjorndonald/lasgcce/internal/handlers"
	"github.com/bjorndonald/lasgcce/internal/middleware"
	"github.com/bjorndonald/lasgcce/internal/validators"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, d *bootstrap.AppDependencies) {
	userRouter := router.Group("user")

	handler := handlers.NewUserHandler(d)

	userRouter.GET("/profile", middleware.JWTMiddleware(d.DatabaseService), handler.UserProfile)
	userRouter.PUT("/", middleware.JWTMiddleware(d.DatabaseService), validators.ValidateUpdateUserProfile, handler.UpdateUserProfile)
}
