package property

import (
	"Praiseson6065/Hypergro-assign/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProperty() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		propertyID := ctx.Param("id")
		if propertyID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Property ID is required",
			})
			return
		}

		property, err := database.GetPropertyByID(ctx, propertyID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"property": property,
		})
	}
}
