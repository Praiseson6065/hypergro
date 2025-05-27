package database

import (
	"Praiseson6065/Hypergro-assign/models"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FirstOrCreateUser(ctx *gin.Context, user *models.User) (string, error) {
	db := GetMongoDB()
	collection := db.Collection("users")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existingUser := &models.User{}
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(existingUser)

	if err != nil && err.Error() == "mongo: no documents in result" {
		user.CreatedAt = time.Now()

		insertResult, err := collection.InsertOne(dbCtx, user)
		if err != nil {
			return "", err
		}

		if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
			return oid.Hex(), nil
		}
		return "", errors.New("failed to get inserted user ID")
	} else if err != nil {
		return "", err
	}

	// User already exists
	return existingUser.ID.Hex(), nil
}

func GetPasswordByMail(ctx *gin.Context, email string) (string, string, error) {
	db := GetMongoDB()
	collection := db.Collection("users")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(dbCtx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return "", "", errors.New("user not found")
		}
		return "", "", err
	}

	return user.Password, user.ID.Hex(), nil
}

func GetUserByID(ctx *gin.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	db := GetMongoDB()
	collection := db.Collection("users")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	filter := bson.M{"_id": objectID}
	err = collection.FindOne(dbCtx, filter).Decode(&user)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func AddToFavorites(ctx *gin.Context, userID string, propertyID string) error {
	// Convert string IDs to ObjectIDs
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	propertyOID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return errors.New("invalid property ID format")
	}

	db := GetMongoDB()
	collection := db.Collection("users")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if property is already in favorites
	var user models.User
	err = collection.FindOne(dbCtx, bson.M{
		"_id": userOID,
		"favorites": bson.M{
			"$in": []primitive.ObjectID{propertyOID},
		},
	}).Decode(&user)

	if err == nil {
		// Property already in favorites
		return nil
	} else if err.Error() != "mongo: no documents in result" {
		return err
	}

	// Add property to favorites
	_, err = collection.UpdateOne(
		dbCtx,
		bson.M{"_id": userOID},
		bson.M{"$addToSet": bson.M{"favorites": propertyOID}},
	)

	return err
}

func RemoveFromFavorites(ctx *gin.Context, userID string, propertyID string) error {

	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	propertyOID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return errors.New("invalid property ID format")
	}

	db := GetMongoDB()
	collection := db.Collection("users")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Remove property from favorites
	_, err = collection.UpdateOne(
		dbCtx,
		bson.M{"_id": userOID},
		bson.M{"$pull": bson.M{"favorites": propertyOID}},
	)

	return err
}

// GetUserFavorites retrieves a list of favorite properties for a user
func GetUserFavorites(ctx context.Context, userID string) ([]models.Property, error) {
	db := GetMongoDB()
	usersCollection := db.Collection("users")
	propertiesCollection := db.Collection("properties")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Convert userID to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Find the user to get their favorites
	var user models.User
	err = usersCollection.FindOne(dbCtx, bson.M{"_id": userObjID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	// If no favorites, return empty slice
	if len(user.Favorites) == 0 {
		return []models.Property{}, nil
	}

	// Query for properties that match the favorite IDs
	cursor, err := propertiesCollection.Find(dbCtx, bson.M{"_id": bson.M{"$in": user.Favorites}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(dbCtx)

	var favoriteProperties []models.Property
	if err := cursor.All(dbCtx, &favoriteProperties); err != nil {
		return nil, err
	}

	return favoriteProperties, nil
}

// AddFavoriteProperty adds a property to a user's favorites
func AddFavoriteProperty(ctx context.Context, userID string, propertyID string) error {
	db := GetMongoDB()
	collection := db.Collection("users")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Convert IDs to ObjectIDs
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	propObjID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return errors.New("invalid property ID format")
	}

	// Check if property exists
	propertyCollection := db.Collection("properties")
	count, err := propertyCollection.CountDocuments(dbCtx, bson.M{"_id": propObjID})
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("property not found")
	}

	// Add to favorites only if not already in favorites
	update := bson.M{
		"$addToSet": bson.M{
			"favorites": propObjID,
		},
	}

	result, err := collection.UpdateOne(dbCtx, bson.M{"_id": userObjID}, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		// Check if user exists
		count, err := collection.CountDocuments(dbCtx, bson.M{"_id": userObjID})
		if err != nil {
			return err
		}
		if count == 0 {
			return errors.New("user not found")
		}
		// If user exists but nothing was modified, property is already in favorites
		return errors.New("property already in favorites")
	}

	return nil
}

// RemoveFavoriteProperty removes a property from a user's favorites
func RemoveFavoriteProperty(ctx context.Context, userID string, propertyID string) error {
	db := GetMongoDB()
	collection := db.Collection("users")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Convert IDs to ObjectIDs
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	propObjID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return errors.New("invalid property ID format")
	}

	// Remove from favorites
	update := bson.M{
		"$pull": bson.M{
			"favorites": propObjID,
		},
	}

	result, err := collection.UpdateOne(dbCtx, bson.M{"_id": userObjID}, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		// Check if user exists
		count, err := collection.CountDocuments(dbCtx, bson.M{"_id": userObjID})
		if err != nil {
			return err
		}
		if count == 0 {
			return errors.New("user not found")
		}
		// If user exists but nothing was modified, property was not in favorites
		return errors.New("property not in favorites")
	}

	return nil
}

// Functions for recommendations

// RecommendProperty adds a property recommendation to a user
func RecommendProperty(ctx context.Context, fromUserID, toUserID, propertyID string) error {
	db := GetMongoDB()
	usersCollection := db.Collection("users")
	propertiesCollection := db.Collection("properties")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Convert IDs to ObjectIDs
	fromUserObjID, err := primitive.ObjectIDFromHex(fromUserID)
	if err != nil {
		return errors.New("invalid recommender user ID format")
	}

	toUserObjID, err := primitive.ObjectIDFromHex(toUserID)
	if err != nil {
		return errors.New("invalid recipient user ID format")
	}

	propObjID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return errors.New("invalid property ID format")
	}

	// Check if property exists
	count, err := propertiesCollection.CountDocuments(dbCtx, bson.M{"_id": propObjID})
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("property not found")
	}

	// Check if recipient user exists
	count, err = usersCollection.CountDocuments(dbCtx, bson.M{"_id": toUserObjID})
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("recipient user not found")
	}

	// Create recommendation
	recommendation := models.Recommendation{
		PropertyID:    propObjID,
		RecommendedBy: fromUserObjID,
		RecommendedAt: time.Now(),
	}

	// Add recommendation to recipient's recommendations
	update := bson.M{
		"$push": bson.M{
			"recommendationsReceived": recommendation,
		},
	}

	_, err = usersCollection.UpdateOne(dbCtx, bson.M{"_id": toUserObjID}, update)
	if err != nil {
		return err
	}

	return nil
}

// GetReceivedRecommendations retrieves recommendations received by a user
func GetReceivedRecommendations(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	db := GetMongoDB()
	usersCollection := db.Collection("users")
	propertiesCollection := db.Collection("properties")
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Convert userID to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Find the user to get their recommendations
	var user models.User
	err = usersCollection.FindOne(dbCtx, bson.M{"_id": userObjID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	// If no recommendations, return empty slice
	if len(user.RecommendationsReceived) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Gather all property IDs and recommender IDs for lookup
	var propertyIDs []primitive.ObjectID
	var recommenderIDs []primitive.ObjectID
	for _, rec := range user.RecommendationsReceived {
		propertyIDs = append(propertyIDs, rec.PropertyID)
		recommenderIDs = append(recommenderIDs, rec.RecommendedBy)
	}

	// Get properties
	propertiesCursor, err := propertiesCollection.Find(dbCtx, bson.M{"_id": bson.M{"$in": propertyIDs}})
	if err != nil {
		return nil, err
	}
	defer propertiesCursor.Close(dbCtx)

	var properties []models.Property
	if err := propertiesCursor.All(dbCtx, &properties); err != nil {
		return nil, err
	}

	// Get recommenders
	recommendersCursor, err := usersCollection.Find(dbCtx, bson.M{"_id": bson.M{"$in": recommenderIDs}})
	if err != nil {
		return nil, err
	}
	defer recommendersCursor.Close(dbCtx)

	var recommenders []models.User
	if err := recommendersCursor.All(dbCtx, &recommenders); err != nil {
		return nil, err
	}

	// Create property and recommender maps for easy lookup
	propertyMap := make(map[string]models.Property)
	for _, property := range properties {
		propertyMap[property.ID.Hex()] = property
	}

	recommenderMap := make(map[string]models.User)
	for _, recommender := range recommenders {
		recommenderMap[recommender.ID.Hex()] = recommender
	}

	// Build detailed recommendations
	var detailedRecommendations []map[string]interface{}
	for _, rec := range user.RecommendationsReceived {
		property, propertyExists := propertyMap[rec.PropertyID.Hex()]
		recommender, recommenderExists := recommenderMap[rec.RecommendedBy.Hex()]

		if propertyExists && recommenderExists {
			detailedRecommendation := map[string]interface{}{
				"property":      property,
				"recommendedBy": map[string]interface{}{"id": recommender.ID, "name": recommender.Name, "email": recommender.Email},
				"recommendedAt": rec.RecommendedAt,
			}
			detailedRecommendations = append(detailedRecommendations, detailedRecommendation)
		}
	}

	return detailedRecommendations, nil
}
