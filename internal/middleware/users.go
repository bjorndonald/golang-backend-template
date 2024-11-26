package middleware

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/bjorndonald/golang-backend-template/constants"
	"github.com/bjorndonald/golang-backend-template/internal/helpers"
	"github.com/bjorndonald/golang-backend-template/internal/models"
	"github.com/bjorndonald/golang-backend-template/internal/repository"
	"github.com/bjorndonald/golang-backend-template/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"strings"
)

type AppError struct {
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
func NewError(message string) *AppError {
	return &AppError{
		Message: message,
	}
}

var (
	constant = constants.New()
)

func OnlyAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		claimsData, ok := claims.(*helpers.AuthTokenJwtClaim)
		if !ok {
			helpers.ReturnJSON(c, "Could not read access token", errors.New("Invalid access token"), http.StatusUnauthorized)
			c.Abort()

			return
		}

		user, _, err := repository.NewUserRepository(db).FindByCondition("id", claimsData.UserId)
		if err != nil {
			helpers.ReturnError(c, "Could not retrieve possible admin user", err, http.StatusUnauthorized)
			c.Abort()
			return
		}
		if user.Role != models.AdminRole {
			helpers.ReturnJSON(c, "Unauthorized access to resource", errors.New("User does not have acces to this request. Please contact admin."), http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("role", user.Role)
		c.Next()
	}
}

func JWTMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the Authorization header
		authHeader := c.GetHeader("Authorization")

		apiQuery := c.Query("access_token")

		if apiQuery != "" {
			authHeader = apiQuery
		}

		if authHeader == "" {
			helpers.ReturnJSON(c, "Missing Authorization Header or Access Token", nil, http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Extract the token from the "Bearer <jwt>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			helpers.ReturnJSON(c, "Invalid Authorization Header or Access Token", nil, http.StatusUnauthorized)
			c.Abort()

			return
		}
		log.Println(tokenString)
		// Parse and validate the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &helpers.AuthTokenJwtClaim{}, func(token *jwt.Token) (interface{}, error) {
			// Provide the same JWT secret key used for signing the tokens
			return []byte(constant.JWTSecretKey), nil
		})
		if err != nil || !token.Valid {
			helpers.ReturnError(c, "Expired Authorization or Access Token", err, http.StatusUnauthorized)
			c.Abort()

			return
		}

		// Extract the claims from the token
		claims, ok := token.Claims.(*helpers.AuthTokenJwtClaim)
		if !ok {
			helpers.ReturnJSON(c, "Invalid claims", nil, http.StatusUnauthorized)
			c.Abort()

			return
		}

		// Attach the claims to the request context for further use
		c.Set("claims", claims)

		_, found, err := repository.NewUserRepository(db).FindByCondition("email", claims.Email)

		if err != nil {
			helpers.ReturnError(c, "Something went wrong", err, http.StatusUnauthorized)
			c.Abort()

			return
		}

		if !found {
			helpers.ReturnJSON(c, "Unauthorized access to resource", nil, http.StatusUnauthorized)
			c.Abort()

			return
		}

		// Proceed to the next middleware or route handler
		c.Next()
	}
}

func CloudinaryUploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the file from the request
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
			c.Abort()
			return
		}

		// Validate the file type
		if !isValidImage(file) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
			c.Abort()
			return
		}

		// Upload the file to Cloudinary
		uploadResult, err := utils.UploadFile(c, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			c.Abort()
			return
		}

		c.Set("cloudinary_url", uploadResult)
		c.Next() // Continue to the next middleware/handler
	}
}

func isValidImage(header *multipart.FileHeader) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	return validTypes[header.Header.Get("Content-Type")]
}
