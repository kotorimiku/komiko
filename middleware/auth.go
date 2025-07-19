package middleware

import (
	"komiko/repo"
	"komiko/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(repo *repo.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// method := c.Request.Method

		// if method == http.MethodOptions {
		// 	c.Status(http.StatusOK)
		// 	c.Abort()
		// 	return
		// }

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		userID := utils.ParseToken(authHeader)

		if userID == "" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		id, err := utils.ParamToUint(userID)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		exists, err := repo.UserRepo.Exists("id = ?", id)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		if !exists {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Set("userID", id)
		c.Next()
	}
}
