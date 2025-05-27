package property

import (
	"Praiseson6065/Hypergro-assign/database"
	"Praiseson6065/Hypergro-assign/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteProperty() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		propertyID := ctx.Param("id")

		userID := middleware.GetUserID(ctx)

		err := database.DeleteAProperty(ctx, propertyID, userID)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Property deleted successfully"})
	}
}
