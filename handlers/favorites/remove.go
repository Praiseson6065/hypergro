package favorites

import (
	"Praiseson6065/Hypergro-assign/database"
	"Praiseson6065/Hypergro-assign/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RemoveFavorite handles DELETE /api/users/{userId}/favorites/{propId} requests
func RemoveFavorite() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get userId and propId from URL params
		userId := ctx.Param("userId")
		propId := ctx.Param("propId")

		if userId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID is required",
			})
			return
		}

		if propId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Property ID is required",
			})
			return
		}

		// Verify the authenticated user matches the requested user ID
		authenticatedUserID := middleware.GetUserID(ctx)
		if authenticatedUserID != userId {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "You can only modify your own favorites",
			})
			return
		}

		// Remove property from favorites
		err := database.RemoveFavoriteProperty(ctx, userId, propId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Property removed from favorites",
		})
	}
}
