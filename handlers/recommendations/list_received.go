package recommendations

import (
	"Praiseson6065/Hypergro-assign/database"
	"Praiseson6065/Hypergro-assign/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListReceivedRecommendations handles GET /api/users/{userId}/recommendations/received requests
func ListReceivedRecommendations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get userId from URL param
		userId := ctx.Param("userId")
		if userId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID is required",
			})
			return
		}

		// Verify the authenticated user matches the requested user ID
		authenticatedUserID := middleware.GetUserID(ctx)
		if authenticatedUserID != userId {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "You can only view your own recommendations",
			})
			return
		}

		// Get received recommendations
		recommendations, err := database.GetReceivedRecommendations(ctx, userId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":          "success",
			"count":           len(recommendations),
			"recommendations": recommendations,
		})
	}
}
