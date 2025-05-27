package favorites

import (
	"Praiseson6065/Hypergro-assign/database"
	"Praiseson6065/Hypergro-assign/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddFavorite handles POST /api/users/{userId}/favorites requests
func AddFavorite() gin.HandlerFunc {
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
				"error": "You can only modify your own favorites",
			})
			return
		}

		// Get property ID from request body
		var request struct {
			PropertyID string `json:"propertyId"`
		}
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Property ID is required",
			})
			return
		}

		// Add property to favorites
		err := database.AddFavoriteProperty(ctx, userId, request.PropertyID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Property added to favorites",
		})
	}
}