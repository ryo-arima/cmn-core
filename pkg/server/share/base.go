package share

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
)

// Common interface for middleware layer.
type Common interface {
	ValidateToken(ctx context.Context, tokenString string) (*model.JWTClaims, error)
}

// ForPublic allows requests without authentication.
func ForPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// ForShare requires a valid IdP-issued JWT token.
func ForShare(commonRepo Common) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateJWTToken(c, commonRepo); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "MIDDLEWARE_AUTH_001",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ForInternal requires a valid IdP-issued JWT token.
func ForInternal(commonRepo Common) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateJWTToken(c, commonRepo); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "MIDDLEWARE_AUTH_001",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ForPrivate requires a valid IdP-issued JWT token and the admin role.
func ForPrivate(commonRepo Common) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateJWTToken(c, commonRepo); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "MIDDLEWARE_AUTH_002",
				"message": "Admin authentication required",
			})
			c.Abort()
			return
		}
		claims, ok := getUserFromContext(c)
		if !ok || claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    "MIDDLEWARE_AUTH_003",
				"message": "Admin role required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// validateJWTToken validates the IdP-issued JWT and sets user context.
func validateJWTToken(c *gin.Context, commonRepo Common) error {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return errors.New("authorization header required")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return errors.New("invalid authorization header format")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return errors.New("token required")
	}

	claims, err := commonRepo.ValidateToken(c.Request.Context(), tokenString)
	if err != nil {
		return err
	}
	setUserContext(c, claims)
	return nil
}

// setUserContext stores user claims in gin context.
func setUserContext(c *gin.Context, claims *model.JWTClaims) {
	c.Set("user_uuid", claims.UUID)
	c.Set("user_email", claims.Email)
	c.Set("user_name", claims.Name)
	c.Set("user_role", claims.Role)
	c.Set("user_claims", claims)
}

// getUserFromContext retrieves user claims from gin context.
func getUserFromContext(c *gin.Context) (*model.JWTClaims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}
	userClaims, ok := claims.(*model.JWTClaims)
	return userClaims, ok
}

// GetUserUUID returns the user UUID from gin context.
func GetUserUUID(c *gin.Context) (string, bool) {
	v, exists := c.Get("user_uuid")
	if !exists {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// GetUserEmail returns the user email from gin context.
func GetUserEmail(c *gin.Context) (string, bool) {
	v, exists := c.Get("user_email")
	if !exists {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// GetUserName returns the user name from gin context.
func GetUserName(c *gin.Context) (string, bool) {
	v, exists := c.Get("user_name")
	if !exists {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// GetUserRole returns the user role from gin context.
func GetUserRole(c *gin.Context) (string, bool) {
	v, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// GetUserClaims returns full user claims from gin context.
func GetUserClaims(c *gin.Context) (*model.JWTClaims, bool) {
	return getUserFromContext(c)
}
