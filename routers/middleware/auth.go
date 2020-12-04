package middleware

import (
	"net/http"
	"strings"

	"cybervein.org/CyberveinDB/models"
	"cybervein.org/CyberveinDB/models/code"
	"cybervein.org/CyberveinDB/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ginMsg := models.GinMsg{C: c}

		respCode := code.CodeTypeOK
		token := c.GetHeader("token")
		address := c.ClientIP()
		var claims *utils.Claims
		var err error
		if token == "" {
			respCode = code.CodeTypePermissionDenied
		} else {
			claims, err = utils.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					respCode = code.CodeTypeTokenTimeoutError
				default:
					respCode = code.CodeTypeTokenInvalidError
				}
			}
		}

		if respCode != code.CodeTypeOK || strings.EqualFold(claims.Address, address) {
			ginMsg.CommonResponse(http.StatusUnauthorized, uint32(respCode), "Invalid Token")
			c.Abort()
			return
		}

		c.Next()
	}
}
