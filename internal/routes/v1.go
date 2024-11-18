package routes

import (
	"github.com/bjorndonald/lasgcce/internal/bootstrap"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup, d *bootstrap.AppDependencies) {

	RegisterUserRoutes(r, d)
	RegisterAuthRoutes(r, d)

}
