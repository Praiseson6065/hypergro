package database

import (
	_ "Praiseson6065/Hypergro-assign/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConfig struct {
	MongoURI      string
	MongoDBName   string
	MongoTimeout  int
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	RedisTimeout  int
}

var env = viper.GetString("ENVIRONMENT")
var MongoClient *mongo.Client

var RedisClient *redis.Client

var dbConfig DBConfig

func InitDB(cfg DBConfig) error {
	const maxRetries = 5
	const retryDelay = 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		if err := connectMongoDB(cfg); err == nil {
			break
		} else if i < maxRetries-1 {
			log.Printf("MongoDB connection attempt %d failed: %v. Retrying in %v...", i+1, err, retryDelay)
			time.Sleep(retryDelay)
			continue
		} else {
			return fmt.Errorf("failed to connect to MongoDB after %d attempts: %w", maxRetries, err)
		}
	}

	for i := 0; i < maxRetries; i++ {
		if err := connectRedis(cfg); err == nil {
			break
		} else if i < maxRetries-1 {
			log.Printf("Redis connection attempt %d failed: %v. Retrying in %v...", i+1, err, retryDelay)
			time.Sleep(retryDelay)
			continue
		} else {
			return fmt.Errorf("failed to connect to Redis after %d attempts: %w", maxRetries, err)
		}
	}

	return nil
}

func connectMongoDB(cfg DBConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.MongoTimeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	MongoClient = client
	log.Println("Connected to MongoDB successfully")
	return nil
}

func connectRedis(cfg DBConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:        cfg.RedisAddr,
		Password:    cfg.RedisPassword,
		DB:          cfg.RedisDB,
		DialTimeout: time.Duration(cfg.RedisTimeout) * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.RedisTimeout)*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Connected to Redis successfully")
	return nil
}

func GetMongoDB() *mongo.Database {

	return MongoClient.Database(viper.GetString(env + ".mongodb.database"))
}

func CloseDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
		log.Println("Disconnected from MongoDB")
	}

	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
		log.Println("Disconnected from Redis")
	}
}

func loadConfig() (DBConfig, error) {

	return DBConfig{
		MongoURI:      viper.GetString(env + ".mongodb.uri"),
		MongoDBName:   viper.GetString(env + ".mongodb.database"),
		MongoTimeout:  viper.GetInt(env + ".mongodb.timeout"),
		RedisAddr:     viper.GetString(env + ".redis.addr"),
		RedisPassword: viper.GetString(env + ".redis.password"),
		RedisDB:       viper.GetInt(env + ".redis.db"),
		RedisTimeout:  viper.GetInt(env + ".redis.timeout"),
	}, nil
}

func init() {
	dbConfig, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to fetch config")
		panic(err)
	}

	err = InitDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect Database")
		panic(err)
	}

}
