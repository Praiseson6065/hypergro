package auth

import (
	"Praiseson6065/Hypergro-assign/database"
	"Praiseson6065/Hypergro-assign/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UserLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var loginRequest LoginRequest
		if err := ctx.ShouldBindBodyWithJSON(&loginRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		hashedPwd, userId, err := database.GetPasswordByMail(ctx, loginRequest.Email)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if !comparePasswords(hashedPwd, loginRequest.Password) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid password",
			})
			return
		}
		token, err := middleware.GenerateToken(userId)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
