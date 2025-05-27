package property

import (
	"Praiseson6065/Hypergro-assign/database"
	"Praiseson6065/Hypergro-assign/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProperty() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userId")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		var propertyRequest models.Property
		if err := ctx.ShouldBindJSON(&propertyRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		userObjID, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID format",
			})
			return
		}

		propertyRequest.CreatedBy = userObjID
		propertyRequest.CreatedAt = time.Now()

		createdProperty, err := database.CreateAProperty(ctx, &propertyRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Property created successfully",
			"data":    createdProperty,
		})
	}
}
