package property

import (
	"Praiseson6065/Hypergro-assign/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserProperties() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		properties, err := database.GetAllPropertiesByUser(ctx)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"properties": properties,
		})
	}
}
