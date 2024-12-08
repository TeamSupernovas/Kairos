package repository

import (
	"context"
	"errors"
	"geodishdiscoveryservice/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DishRepository struct {
	collection *mongo.Collection
}

// NewDishRepository initializes a new MongoDB DishRepository
func NewDishRepository(client *mongo.Client, dbName, collectionName string) *DishRepository {
	return &DishRepository{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

// CreateDish inserts a new dish into the collection
func (r *DishRepository) CreateDish(ctx context.Context, dish *model.Dish) error {
	dish.CreatedAt = time.Now()
	dish.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, dish)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDish updates an existing dish by dishID
func (r *DishRepository) UpdateDish(ctx context.Context, dishID string, dish *model.Dish) error {
	filter := bson.M{"DishID": dishID}
	dish.UpdatedAt = time.Now()
	update := bson.M{
		"$set": dish,
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no dish found with the provided ID")
	}
	return nil
}

// GetDishByID retrieves a dish by its ID
func (r *DishRepository) GetDishByID(ctx context.Context, dishID string) (*model.Dish, error) {
	filter := bson.M{"DishID": dishID}
	var dish model.Dish
	err := r.collection.FindOne(ctx, filter).Decode(&dish)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("no dish found with the provided ID")
	} else if err != nil {
		return nil, err
	}
	return &dish, nil
}

// DeleteDishByID soft deletes a dish by setting its deleted_at field
func (r *DishRepository) DeleteDishByID(ctx context.Context, dishID string) error {
	filter := bson.M{"DishID": dishID}
	update := bson.M{
		"$set": bson.M{"DeletedAt": time.Now()},
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no dish found with the provided ID")
	}
	return nil
}

// GetDishesByLocation retrieves dishes within a radius of a given latitude and longitude
func (r *DishRepository) GetDishesByLocation(
    ctx context.Context,
    lat, long, radius float64,
    mealCourse, dietaryCategory string,
    maxPrice float64,
    skip, limit int,
) ([]model.Dish, int, error) {
    // Convert radius from meters to radians (required for $geoWithin)
    radiusInRadians := radius / 6378100.0 // Earth's radius in meters

    // Build the base filter for geospatial query
    filter := bson.M{
        "location": bson.M{
            "$geoWithin": bson.M{
                "$centerSphere": []interface{}{
                    []float64{long, lat}, // [longitude, latitude]
                    radiusInRadians,     // Radius in radians
                },
            },
        },
    }

    // Add optional filters
    if mealCourse != "" {
        filter["MealCourse"] = mealCourse
    }
    if dietaryCategory != "" {
        filter["DietaryCategory"] = dietaryCategory
    }
    if maxPrice > 0 {
        filter["Price"] = bson.M{"$lte": maxPrice}
    }

    // Count the total number of documents matching the filter
    totalCount, err := r.collection.CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, err
    }

    // Apply pagination options
    findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

    // Query the collection
    cursor, err := r.collection.Find(ctx, filter, findOptions)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)

    // Decode documents into a slice of dishes
    var dishes []model.Dish
    for cursor.Next(ctx) {
        var dish model.Dish
        if err := cursor.Decode(&dish); err != nil {
            return nil, 0, err
        }
        dishes = append(dishes, dish)
    }
    if err := cursor.Err(); err != nil {
        return nil, 0, err
    }

    return dishes, int(totalCount), nil
}
