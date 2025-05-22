package middlewares

import (
	"github.com/Teemo4621/Hospital-Api/configs"
	"github.com/Teemo4621/Hospital-Api/pkgs/utils"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	cfg *configs.Config
}

func NewAuthMiddleware(cfg *configs.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

func (a *AuthMiddleware) JwtAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			utils.UnauthorizedResponse(c, "Unauthorized")
			c.Abort()
			return
		}

		tokenData, err := utils.ParseAccessToken(a.cfg, tokenString)
		if err != nil {
			utils.UnauthorizedResponse(c, "Unauthorized")
			c.Abort()
			return
		}

		c.Set("user_data", tokenData)
		c.Next()
	}
}
