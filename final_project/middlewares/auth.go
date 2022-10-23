package middlewares

import (
	"final_project/dto"
	"final_project/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyToken, err := helpers.VerifyToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Error: "please login first",
			})
			return
		}
		c.Set("userData", verifyToken)
		c.Next()
	}
}
