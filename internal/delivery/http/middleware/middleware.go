package http

import (
	"payroll-system/internal/delivery/dto"
	httpError "payroll-system/internal/error_const"
	"payroll-system/internal/utils"

	"github.com/gin-gonic/gin"
)

func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, dto.NewErrorResponse("Unauthorized", httpError.ErrJWTTokenRequired))
			c.Abort()
			return
		}

		valid, err := utils.IsValidJWT(utils.GetTokenWithoutBearer(token))
		if !valid || err != nil {
			c.JSON(401, dto.NewErrorResponse("Unauthorized", httpError.ErrJWTTokenInvalid))
			c.Abort()
			return
		}

		c.Next()
	}
}

func CheckRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, dto.NewErrorResponse("Unauthorized", httpError.ErrJWTTokenRequired))
			c.Abort()
			return
		}
		claims, err := utils.GetClaimsFromJWT(utils.GetTokenWithoutBearer(token))
		if err != nil {
			c.JSON(401, dto.NewErrorResponse("Unauthorized", httpError.ErrJWTTokenInvalid))
			c.Abort()
			return
		}
		role := claims.Role
		if role != requiredRole {
			c.JSON(403, dto.NewErrorResponse("Forbidden", httpError.ErrNotAllowedAccess))
			c.Abort()
			return
		}
		c.Next()
	}
}
