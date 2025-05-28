package property

import (
	"Praiseson6065/Hypergro-assign/database"
	"Praiseson6065/Hypergro-assign/middleware"
	"Praiseson6065/Hypergro-assign/models"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ImportPropertiesFromCSV handles the import of properties from a CSV file
func ImportPropertiesFromCSV() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := middleware.GetUserID(ctx)

		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		file, err := ctx.FormFile("properties_csv")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded or invalid form field"})
			return
		}

		if !strings.HasSuffix(file.Filename, ".csv") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed"})
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open the uploaded file"})
			return
		}
		defer openedFile.Close()

		reader := csv.NewReader(openedFile)
		records, err := reader.ReadAll()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse CSV file: " + err.Error()})
			return
		}

		if len(records) < 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "CSV file has insufficient data"})
			return
		}

		headers := records[0]

		headerMap := make(map[string]int)
		for i, header := range headers {
			headerMap[strings.TrimSpace(strings.ToLower(header))] = i
		}

		requiredFields := []string{"title", "type", "price", "state", "city"}
		for _, field := range requiredFields {
			if _, exists := headerMap[field]; !exists {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("CSV is missing required field: %s", field),
				})
				return
			}
		}

		var createdProperties []models.Property
		var errorMessages []string

		for i, record := range records[1:] {

			property := models.Property{
				ID:         primitive.NewObjectID(),
				CreatedBy:  userObjID,
				CreatedAt:  time.Now(),
				IsVerified: false,
			}

			rowNum := i + 2

			if idx, exists := headerMap["title"]; exists && idx < len(record) {
				property.Title = strings.TrimSpace(record[idx])
				if property.Title == "" {
					errorMessages = append(errorMessages, fmt.Sprintf("Row %d: Title is required", rowNum))
					continue
				}
			}

			if idx, exists := headerMap["type"]; exists && idx < len(record) {
				property.Type = strings.TrimSpace(record[idx])
				if property.Type == "" {
					errorMessages = append(errorMessages, fmt.Sprintf("Row %d: Type is required", rowNum))
					continue
				}
			}

			if idx, exists := headerMap["price"]; exists && idx < len(record) {
				priceStr := strings.TrimSpace(record[idx])
				price, err := strconv.ParseInt(priceStr, 10, 64)
				if err != nil || price <= 0 {
					errorMessages = append(errorMessages, fmt.Sprintf("Row %d: Invalid price value", rowNum))
					continue
				}
				property.Price = price
			}

			if idx, exists := headerMap["state"]; exists && idx < len(record) {
				property.State = strings.TrimSpace(record[idx])
				if property.State == "" {
					errorMessages = append(errorMessages, fmt.Sprintf("Row %d: State is required", rowNum))
					continue
				}
			}

			if idx, exists := headerMap["city"]; exists && idx < len(record) {
				property.City = strings.TrimSpace(record[idx])
				if property.City == "" {
					errorMessages = append(errorMessages, fmt.Sprintf("Row %d: City is required", rowNum))
					continue
				}
			}

			if idx, exists := headerMap["areasqft"]; exists && idx < len(record) && strings.TrimSpace(record[idx]) != "" {
				areaSqFt, err := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64)
				if err == nil {
					property.AreaSqFt = areaSqFt
				}
			}

			if idx, exists := headerMap["bedrooms"]; exists && idx < len(record) && strings.TrimSpace(record[idx]) != "" {
				bedrooms, err := strconv.Atoi(strings.TrimSpace(record[idx]))
				if err == nil {
					property.Bedrooms = bedrooms
				}
			}

			if idx, exists := headerMap["bathrooms"]; exists && idx < len(record) && strings.TrimSpace(record[idx]) != "" {
				bathrooms, err := strconv.Atoi(strings.TrimSpace(record[idx]))
				if err == nil {
					property.Bathrooms = bathrooms
				}
			}

			if idx, exists := headerMap["amenities"]; exists && idx < len(record) && strings.TrimSpace(record[idx]) != "" {
				amenitiesStr := strings.TrimSpace(record[idx])
				if amenitiesStr != "" {
					property.Amenities = strings.Split(amenitiesStr, "|")

					for i, amenity := range property.Amenities {
						property.Amenities[i] = strings.TrimSpace(amenity)
					}
				}
			}

			if idx, exists := headerMap["furnished"]; exists && idx < len(record) {
				property.Furnished = strings.TrimSpace(record[idx])
			}

			if idx, exists := headerMap["availablefrom"]; exists && idx < len(record) && strings.TrimSpace(record[idx]) != "" {
				dateStr := strings.TrimSpace(record[idx])
				date, err := time.Parse("2006-01-02", dateStr)
				if err == nil {
					property.AvailableFrom = date
				}
			}

			if idx, exists := headerMap["listedby"]; exists && idx < len(record) {
				property.ListedBy = strings.TrimSpace(record[idx])
			}

			if idx, exists := headerMap["tags"]; exists && idx < len(record) && strings.TrimSpace(record[idx]) != "" {
				tagsStr := strings.TrimSpace(record[idx])
				if tagsStr != "" {
					property.Tags = strings.Split(tagsStr, "|")

					for i, tag := range property.Tags {
						property.Tags[i] = strings.TrimSpace(tag)
					}
				}
			}

			if idx, exists := headerMap["colortheme"]; exists && idx < len(record) {
				property.ColorTheme = strings.TrimSpace(record[idx])
			}

			if idx, exists := headerMap["rating"]; exists && idx < len(record) && strings.TrimSpace(record[idx]) != "" {
				rating, err := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64)
				if err == nil {
					property.Rating = rating
				}
			}

			if idx, exists := headerMap["isverified"]; exists && idx < len(record) && strings.TrimSpace(record[idx]) != "" {
				isVerifiedStr := strings.ToLower(strings.TrimSpace(record[idx]))
				if isVerifiedStr == "true" || isVerifiedStr == "yes" || isVerifiedStr == "1" {
					property.IsVerified = true
				}
			}

			if idx, exists := headerMap["listingtype"]; exists && idx < len(record) {
				property.ListingType = strings.TrimSpace(record[idx])
			}

			createdProperty, err := database.CreateAProperty(ctx, &property)
			if err != nil {
				errorMessages = append(errorMessages, fmt.Sprintf("Row %d: Failed to create property: %s", rowNum, err.Error()))
				continue
			}

			createdProperties = append(createdProperties, *createdProperty)
		}

		if len(errorMessages) > 0 && len(createdProperties) == 0 {

			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Failed to import any properties",
				"details": errorMessages,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":            "Properties imported successfully",
			"total_processed":    len(records) - 1,
			"total_created":      len(createdProperties),
			"created_properties": createdProperties,
			"errors":             errorMessages,
		})
	}
}
