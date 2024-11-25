package routes

import (
	"github.com/bjorndonald/golang-backend-template/internal/bootstrap"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup, d *bootstrap.AppDependencies) {

	RegisterUserRoutes(r, d)
	RegisterAuthRoutes(r, d)

}
