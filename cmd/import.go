package main

import (
	"Praiseson6065/Hypergro-assign/database"
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Import() {
	f, err := os.Open("db424fd9fb74_1748258398689.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	headers, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	db := database.GetMongoDB()
	coll := db.Collection("properties")

	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "city", Value: 1}}},
		{Keys: bson.D{{Key: "price", Value: 1}}},
		{Keys: bson.D{{Key: "state", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
	})

	var docs []interface{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		doc := bson.M{}
		doc["_id"] = primitive.NewObjectID()

		for i, v := range record {
			key := headers[i]

			if key == "propId" || key == "id" {
				continue
			}

			switch key {
			case "price", "areaSqFt":

				if val, err := strconv.ParseInt(v, 10, 64); err == nil {
					doc[key] = val
				} else {
					doc[key] = 0
				}
			case "bedrooms", "bathrooms":

				if val, err := strconv.Atoi(v); err == nil {
					doc[key] = val
				} else {
					doc[key] = 0
				}
			case "rating":

				if val, err := strconv.ParseFloat(v, 64); err == nil {
					doc[key] = val
				} else {
					doc[key] = 0.0
				}
			case "isVerified":

				doc[key] = v == "True"
			case "amenities", "tags":

				doc[key] = strings.Split(v, "|")
			case "availableFrom":

				if t, err := time.Parse("2006-01-02", v); err == nil {
					doc[key] = t
				} else {
					doc[key] = time.Now()
				}
			default:
				doc[key] = v
			}
		}
		doc["createdAt"] = time.Now()

		objID, _ := primitive.ObjectIDFromHex("6835c855bdcc74cfb350e6c4")
		doc["createdBy"] = objID

		docs = append(docs, doc)

		if len(docs) >= 1000 {
			if _, err := coll.InsertMany(ctx, docs); err != nil {
				log.Fatal(err)
			}
			log.Printf("Inserted %d documents", len(docs))
			docs = docs[:0]
		}
	}

	if len(docs) > 0 {
		if _, err := coll.InsertMany(ctx, docs); err != nil {
			log.Fatal(err)
		}
		log.Printf("Inserted %d documents", len(docs))
	}

	log.Println("Import complete!")
}
