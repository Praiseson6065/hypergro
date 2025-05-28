package database

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ShortTerm  = 5 * time.Minute
	MediumTerm = 1 * time.Hour
	LongTerm   = 24 * time.Hour

	PropertyKeyPrefix      = "property:"
	PropertiesKeyPrefix    = "properties:"
	UserFavoritesKeyPrefix = "user:favorites:"
	UserKeyPrefix          = "user:"
)

func GetFromCache(ctx *gin.Context, key string, result interface{}) (bool, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return false, nil
	}

	err = json.Unmarshal([]byte(val), result)
	if err != nil {
		return false, err
	}

	return true, nil
}

func SetInCache(ctx *gin.Context, key string, value interface{}, expiry time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, key, data, expiry).Err()
}

func DeleteFromCache(ctx *gin.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

func DeleteByPattern(ctx *gin.Context, pattern string) error {
	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

func ClearPropertyCache(ctx *gin.Context, propertyID string) {
	key := PropertyKeyPrefix + propertyID
	err := DeleteFromCache(ctx, key)
	if err != nil {
		log.Printf("Error clearing property cache for %s: %v", propertyID, err)
	}

	// Also clear any property lists that might contain this property
	err = DeleteByPattern(ctx, PropertiesKeyPrefix+"*")
	if err != nil {
		log.Printf("Error clearing properties list cache: %v", err)
	}
}

func ClearUserFavoritesCache(ctx *gin.Context, userID string) {
	key := UserFavoritesKeyPrefix + userID
	err := DeleteFromCache(ctx, key)
	if err != nil {
		log.Printf("Error clearing user favorites cache for %s: %v", userID, err)
	}
}

func ClearUserCache(ctx *gin.Context, userID string) {
	key := UserKeyPrefix + userID
	err := DeleteFromCache(ctx, key)
	if err != nil {
		log.Printf("Error clearing user cache for %s: %v", userID, err)
	}
}
