package property

import (
	"Praiseson6065/Hypergro-assign/database"
	"Praiseson6065/Hypergro-assign/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateProperty() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		propertyID := ctx.Param("id")

		userID := middleware.GetUserID(ctx)
		if userID == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		
		var updateData map[string]interface{}
		if err := ctx.ShouldBindJSON(&updateData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		objID, err := primitive.ObjectIDFromHex(propertyID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID format"})
			return
		}


		updatedProperty, err := database.UpdateAProperty(ctx, updateData, objID)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Property updated successfully", "property": updatedProperty})
	}
}
