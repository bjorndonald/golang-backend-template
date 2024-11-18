package validators

import (
	"github.com/bjorndonald/lasgcce/internal/handlers"
	"github.com/gin-gonic/gin"
)

func ValidatePageInputSchema(c *gin.Context) {
	var body handlers.PageInput
	bindAndValidate(c, &body)
	c.Set("validatedRequestBody", body)
	c.Next()
}
