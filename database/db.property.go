package database

import (
	"Praiseson6065/Hypergro-assign/middleware"
	"Praiseson6065/Hypergro-assign/models"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllPropertiesByUser(ctx *gin.Context) ([]models.Property, error) {
	userID := middleware.GetUserID(ctx)
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

	var properties []models.Property
	if err := cursor.All(dbCtx, &properties); err != nil {
		return nil, err
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

	return nil
}

func GetPropertyByID(ctx context.Context, propertyID string) (*models.Property, error) {
	db := GetMongoDB()
	collection := db.Collection("properties")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	propObjID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return nil, errors.New("invalid property ID format")
	}

	filter := bson.M{"_id": propObjID}

	var property models.Property
	err = collection.FindOne(dbCtx, filter).Decode(&property)
	if err != nil {
		return nil, err
	}

	return &property, nil
}

func GetAllProperties(ctx context.Context, filters bson.M) ([]models.Property, error) {
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

	return properties, nil
}
