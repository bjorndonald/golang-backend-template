package routes

import (
	"github.com/bjorndonald/lasgcce/internal/bootstrap"
	"github.com/bjorndonald/lasgcce/internal/handlers"
	"github.com/bjorndonald/lasgcce/internal/middleware"
	"github.com/bjorndonald/lasgcce/internal/validators"

	"github.com/gin-gonic/gin"
)

func RegisterWorkRoutes(router *gin.RouterGroup, d *bootstrap.AppDependencies) {
	workRouter := router.Group("work")

	handler := handlers.NewWorkHandler(d)

	workRouter.POST("/", middleware.JWTMiddleware(d.DatabaseService), middleware.OnlyAdmin(d.DatabaseService), validators.ValidatePageInputSchema, handler.CreatePage)
	workRouter.PUT("/", middleware.JWTMiddleware(d.DatabaseService), middleware.OnlyAdmin(d.DatabaseService), validators.ValidatePageInputSchema, handler.UpdatePage)
}
