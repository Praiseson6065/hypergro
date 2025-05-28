package database

import (
	"Praiseson6065/Hypergro-assign/middleware"
	"Praiseson6065/Hypergro-assign/models"
	"context"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPropertyByID(ctx *gin.Context, propertyID string) (*models.Property, error) {

	cacheKey := PropertyKeyPrefix + propertyID
	var property models.Property
	found, err := GetFromCache(ctx, cacheKey, &property)
	if err != nil {
		log.Printf("Error retrieving property from cache: %v", err)
	}

	if found {
		return &property, nil
	}

	// Not in cache, get from database
	db := GetMongoDB()
	collection := db.Collection("properties")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	propObjID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return nil, errors.New("invalid property ID format")
	}

	filter := bson.M{"_id": propObjID}

	err = collection.FindOne(dbCtx, filter).Decode(&property)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	err = SetInCache(ctx, cacheKey, property, MediumTerm)
	if err != nil {
		log.Printf("Error caching property: %v", err)
	}

	return &property, nil
}

func GetAllProperties(ctx *gin.Context, filters bson.M) ([]models.Property, error) {
	// Create a cache key based on the filters
	filterBytes, err := bson.Marshal(filters)
	if err != nil {
		log.Printf("Error marshaling filters for cache key: %v", err)
		// Continue without caching
	} else {
		// Use the filters in the cache key to ensure uniqueness
		cacheKey := PropertiesKeyPrefix + string(filterBytes)
		var properties []models.Property
		found, err := GetFromCache(ctx, cacheKey, &properties)
		if err != nil {
			log.Printf("Error retrieving properties from cache: %v", err)
		}

		if found {
			return properties, nil
		}
	}

	// Not in cache or error creating cache key, get from database
	db := GetMongoDB()
	collection := db.Collection("properties")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if filters == nil {
		filters = bson.M{}
	}

	cursor, err := collection.Find(dbCtx, filters)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(dbCtx)

	var properties []models.Property
	if err := cursor.All(dbCtx, &properties); err != nil {
		return nil, err
	}

	// Store in cache for future requests if we successfully created a cache key
	if filterBytes != nil {
		cacheKey := PropertiesKeyPrefix + string(filterBytes)
		err = SetInCache(ctx, cacheKey, properties, ShortTerm) // Use a shorter cache time for lists
		if err != nil {
			log.Printf("Error caching properties: %v", err)
		}
	}

	return properties, nil
}

func GetAllPropertiesByUser(ctx *gin.Context) ([]models.Property, error) {
	userID := middleware.GetUserID(ctx)

	// Try to get from cache first
	cacheKey := PropertiesKeyPrefix + "user:" + userID
	var properties []models.Property
	found, err := GetFromCache(ctx, cacheKey, &properties)
	if err != nil {
		log.Printf("Error retrieving user properties from cache: %v", err)
	}

	if found {
		return properties, nil
	}

	// Not in cache, get from database
	db := GetMongoDB()
	collection := db.Collection("properties")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	filter := bson.M{"createdBy": userObjID}

	cursor, err := collection.Find(dbCtx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(dbCtx)

	if err := cursor.All(dbCtx, &properties); err != nil {
		return nil, err
	}

	// Store in cache for future requests
	err = SetInCache(ctx, cacheKey, properties, ShortTerm)
	if err != nil {
		log.Printf("Error caching user properties: %v", err)
	}

	return properties, nil
}

func CreateAProperty(ctx *gin.Context, property *models.Property) (*models.Property, error) {
	db := GetMongoDB()
	collection := db.Collection("properties")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	property.CreatedAt = time.Now()

	_, err := collection.InsertOne(dbCtx, property)
	if err != nil {
		return nil, err
	}

	// Invalidate any cached property lists since we added a new property
	err = DeleteByPattern(ctx, PropertiesKeyPrefix+"*")
	if err != nil {
		log.Printf("Error clearing properties cache: %v", err)
	}

	return property, nil

}

func UpdateAProperty(ctx *gin.Context, propertyUpdates map[string]interface{}, propertyID primitive.ObjectID) (*models.Property, error) {
	userID := middleware.GetUserID(ctx)
	db := GetMongoDB()
	collection := db.Collection("properties")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	filter := bson.M{
		"_id":       propertyID,
		"createdBy": userObjID,
	}

	count, err := collection.CountDocuments(dbCtx, filter)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("property not found or you don't have permission to update it")
	}

	// Only update fields that were provided in the request
	update := bson.M{
		"$set": propertyUpdates,
	}

	_, err = collection.UpdateOne(dbCtx, filter, update)
	if err != nil {
		return nil, err
	}

	// Fetch the updated property to return
	var updatedProperty models.Property
	err = collection.FindOne(dbCtx, bson.M{"_id": propertyID}).Decode(&updatedProperty)
	if err != nil {
		return nil, err
	}

	// Clear the cache for this property
	ClearPropertyCache(ctx, propertyID.Hex())

	return &updatedProperty, nil
}

func DeleteAProperty(ctx *gin.Context, propertyID string, userID string) error {
	db := GetMongoDB()
	collection := db.Collection("properties")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	propObjID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return errors.New("invalid property ID format")
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	filter := bson.M{
		"_id":       propObjID,
		"createdBy": userObjID,
	}

	result, err := collection.DeleteOne(dbCtx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("property not found or you don't have permission to delete it")
	}

	// Clear the cache for this property and any property lists
	ClearPropertyCache(ctx, propertyID)

	return nil
}
