package property

import (
	"Praiseson6065/Hypergro-assign/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ListProperties handles GET /api/properties requests for listing and filtering properties
func ListProperties() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Build filters based on query parameters
		filters := bson.M{}

		// Handle common property filters
		if propertyType := ctx.Query("type"); propertyType != "" {
			filters["type"] = propertyType
		}

		if city := ctx.Query("city"); city != "" {
			filters["city"] = city
		}

		if state := ctx.Query("state"); state != "" {
			filters["state"] = state
		}

		if minPrice := ctx.Query("minPrice"); minPrice != "" {
			price, err := strconv.ParseInt(minPrice, 10, 64)
			if err == nil {
				filters["price"] = bson.M{"$gte": price}
			}
		}

		if maxPrice := ctx.Query("maxPrice"); maxPrice != "" {
			price, err := strconv.ParseInt(maxPrice, 10, 64)
			if err == nil {
				if priceFilter, ok := filters["price"].(bson.M); ok {
					priceFilter["$lte"] = price
				} else {
					filters["price"] = bson.M{"$lte": price}
				}
			}
		}

		if bedrooms := ctx.Query("bedrooms"); bedrooms != "" {
			beds, err := strconv.Atoi(bedrooms)
			if err == nil {
				filters["bedrooms"] = beds
			}
		}

		if bathrooms := ctx.Query("bathrooms"); bathrooms != "" {
			baths, err := strconv.Atoi(bathrooms)
			if err == nil {
				filters["bathrooms"] = baths
			}
		}

		if furnished := ctx.Query("furnished"); furnished != "" {
			filters["furnished"] = furnished
		}

		properties, err := database.GetAllProperties(ctx, filters)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"count":      len(properties),
			"properties": properties,
		})
	}
}
