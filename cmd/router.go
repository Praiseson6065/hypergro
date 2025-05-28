package main

import (
	"Praiseson6065/Hypergro-assign/handlers/auth"
	"Praiseson6065/Hypergro-assign/handlers/favorites"
	"Praiseson6065/Hypergro-assign/handlers/property"
	"Praiseson6065/Hypergro-assign/handlers/recommendations"
	"Praiseson6065/Hypergro-assign/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.Engine) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/signup", auth.UserSignup())
		authRoutes.POST("/login", auth.UserLogin())
	}
}

func ApiRouter(r *gin.Engine) {

	apiRoutes := r.Group("/api")

	propertyRoutes := apiRoutes.Group("/properties")
	{

		propertyRoutes.GET("", property.ListProperties())
		propertyRoutes.GET("/:id", property.GetProperty())

		authenticatedPropertyRoutes := propertyRoutes.Group("")
		authenticatedPropertyRoutes.Use(middleware.Authenicator())
		{
			authenticatedPropertyRoutes.POST("", property.CreateProperty())
			authenticatedPropertyRoutes.PUT("/:id", property.UpdateProperty())
			authenticatedPropertyRoutes.DELETE("/:id", property.DeleteProperty())
			authenticatedPropertyRoutes.POST("/import-csv", property.ImportPropertiesFromCSV())
		}
	}

	userRoutes := apiRoutes.Group("/users")
	userRoutes.Use(middleware.Authenicator())
	{

		userRoutes.GET("/:userId/favorites", favorites.ListUserFavorites())
		userRoutes.POST("/:userId/favorites", favorites.AddFavorite())
		userRoutes.DELETE("/:userId/favorites/:propId", favorites.RemoveFavorite())
		userRoutes.GET("/:userId/recommendations/received", recommendations.ListReceivedRecommendations())
	}

	recommendationRoutes := apiRoutes.Group("/recommendations")
	recommendationRoutes.Use(middleware.Authenicator())
	{
		recommendationRoutes.POST("", recommendations.CreateRecommendation())
	}

}
