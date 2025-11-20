package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/phanindra08/snap-attend/internal/shared/config"
	"github.com/phanindra08/snap-attend/internal/shared/models"
)

func RequestLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()

		method := context.Request.Method
		path := context.Request.URL.Path
		rawQuery := context.Request.URL.RawQuery
		clientIP := context.ClientIP()

		if rawQuery != "" {
			path += "?" + rawQuery
		}

		// Process request
		context.Next()

		latency := time.Since(start)
		statusCode := context.Writer.Status()

		log.Printf(
			"[HTTP] %s | %3d | %v | %15s | %-7s %s",
			time.Now().Format("2006-01-02 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

// CORSMiddleware - helps to handle CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Or restrict to FE domain
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Headers",
			"Authorization, Content-Type, X-Requested-With, X-CSRF-Token, Accept, Origin")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if context.Request.Method == http.MethodOptions {
			context.AbortWithStatus(204)
			return
		}

		context.Next()
	}
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")

		if authHeader == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// Extract token string (removes "Bearer ")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.GetConfig().Server.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Printf("JWT parse error: %v", err)
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Extract user_id + role
		userId, ok := claims["user_id"].(float64)
		if !ok {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid role in token"})
			return
		}

		// Store values in Gin context
		context.Set("userID", uint(userId))
		context.Set("role", userRole)

		context.Next()
	}
}

// GenerateToken - Creates a new JWT token with user ID and role
func GenerateToken(userID uint, role models.UserRoles) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    string(role),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	})

	return token.SignedString(config.GetConfig().Server.JwtSecret)
}

// RequireAdmin - Allows only users with Admins role.
func RequireAdmin() gin.HandlerFunc {
	return RequireRoles(models.Admin)
}

// RequireProfessor - Allows only users with Professors role
func RequireProfessor() gin.HandlerFunc {
	return RequireRoles(models.Professor)
}

// RequireStudent - Allows only users with Students role
func RequireStudent() gin.HandlerFunc {
	return RequireRoles(models.Student)
}

// RequireRoles - ensures that the authenticated user has at least one of the allowed roles.
// Usage: r.Use(JWTAuthMiddleware(), RequireRoles(models.Admin, models.Professor))
func RequireRoles(allowedRoles ...models.UserRoles) gin.HandlerFunc {
	return func(context *gin.Context) {
		roleVal, exists := context.Get("role")
		if !exists {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: role is missing in the context"})
			context.Abort()
			return
		}

		roleStr, ok := roleVal.(string)
		if !ok {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid role type"})
			context.Abort()
			return
		}

		// Check if roleStr matches any of the allowed roles
		for _, allowed := range allowedRoles {
			if roleStr == string(allowed) {
				context.Next()
				return
			}
		}

		context.JSON(http.StatusForbidden, gin.H{"error": "Access denied: insufficient permissions for the user"})
		context.Abort()
	}
}
