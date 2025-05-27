package recommendations

import (
	"Praiseson6065/Hypergro-assign/database"
	"Praiseson6065/Hypergro-assign/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRecommendation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fromUserID := middleware.GetUserID(ctx)
		if fromUserID == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		var request struct {
			ToUserID   string `json:"toUserId"`
			PropertyID string `json:"propertyId"`
		}
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request: " + err.Error(),
			})
			return
		}

		if request.ToUserID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Recipient user ID is required",
			})
			return
		}

		if request.PropertyID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Property ID is required",
			})
			return
		}

		err := database.RecommendProperty(ctx, fromUserID, request.ToUserID, request.PropertyID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Property recommendation sent successfully",
		})
	}
}
