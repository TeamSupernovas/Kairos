package service

import (
	"context"
	"encoding/json"
	"fmt"
	"geodishdiscoveryservice/internal/model"
	"geodishdiscoveryservice/internal/repository"
	"time"

	"github.com/segmentio/kafka-go"
)

type DishService struct {
	dishRepo    *repository.DishRepository
	kafkaService *KafkaService
}

func NewDishService(dishRepo *repository.DishRepository, kafkaService *KafkaService) *DishService {
	return &DishService{
		dishRepo:    dishRepo,
		kafkaService: kafkaService,
	}
}

// GetNearbyDishes retrieves dishes near a specified location
func (s *DishService) GetNearbyDishes(
	ctx context.Context,
	latitude, longitude, radius float64,
	mealCourse, dietaryCategory string,
	maxPrice float64,
	page, pageSize int,
) ([]model.Dish, int, error) {
	skip := (page - 1) * pageSize

	dishes, total, err := s.dishRepo.GetDishesByLocation(
		ctx,
		latitude, longitude, radius,
		mealCourse, dietaryCategory, maxPrice,
		skip, pageSize,
	)
	if err != nil {
		return nil, 0, err
	}
	return dishes, total, nil
}

// SubscribeToDishTopics subscribes to Kafka topics and processes messages
func (s *DishService) SubscribeToDishTopics(ctx context.Context) error {
	topics := map[string]func(msg kafka.Message) error{
		"dish-management-service.dish.created": s.handleDishCreated,
		"dish-management-service.dish.updated": s.handleDishUpdated,
		"dish-management-service.dish.deleted": s.handleDishDeleted,
	}

	for topic, handler := range topics {
		if err := s.kafkaService.SubscribeToTopic(topic, handler); err != nil {
			return fmt.Errorf("failed to subscribe to topic %s: %v", topic, err)
		}
	}
	return nil
}

// Handle "Dish Created" event
func (s *DishService) handleDishCreated(msg kafka.Message) error {
	eventData, err := s.parseMessage(msg)
	if err != nil {
		return err
	}

	dish, err := s.constructDishFromEvent(eventData)
	if err != nil {
		return fmt.Errorf("failed to construct dish from event: %v", err)
	}

	return s.dishRepo.CreateDish(context.Background(), dish)
}

// Handle "Dish Updated" event
func (s *DishService) handleDishUpdated(msg kafka.Message) error {
	eventData, err := s.parseMessage(msg)
	if err != nil {
		return err
	}

	dish, err := s.constructDishFromEvent(eventData)
	if err != nil {
		return fmt.Errorf("failed to construct dish from event: %v", err)
	}

	return s.dishRepo.UpdateDish(context.Background(), dish.DishID, dish)
}

// Handle "Dish Deleted" event
func (s *DishService) handleDishDeleted(msg kafka.Message) error {
	eventData, err := s.parseMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to parse Dish Deleted event: %v", err)
	}

	dishID, ok := eventData["dish_id"].(string)
	if !ok || dishID == "" {
		return fmt.Errorf("invalid or missing dish_id in event data")
	}

	return s.dishRepo.DeleteDishByID(context.Background(), dishID)
}

// Helper to parse Kafka message
func (s *DishService) parseMessage(msg kafka.Message) (map[string]interface{}, error) {
	var eventData map[string]interface{}
	if err := json.Unmarshal(msg.Value, &eventData); err != nil {
		return nil, fmt.Errorf("failed to parse Kafka message: %v", err)
	}
	return eventData, nil
}

// Helper to construct Dish model from event data
func (s *DishService) constructDishFromEvent(eventData map[string]interface{}) (*model.Dish, error) {
	availableUntil := parseOptionalTime(eventData, "available_until")

	createdAt, err := parseTime(eventData, "created_at")
	if err != nil {
		return nil, err
	}

	updatedAt, err := parseTime(eventData, "updated_at")
	if err != nil {
		return nil, err
	}

	deletedAt := parseOptionalTime(eventData, "deleted_at")

	geoPoint, err := parseGeoPoint(eventData)
	if err != nil {
		return nil, err
	}

	address, err := parseAddress(eventData)
	if err != nil {
		return nil, err
	}

	dish := &model.Dish{
		DishID:           eventData["dish_id"].(string),
		DishName:         eventData["dish_name"].(string),
		ChefID:           eventData["chef_id"].(string),
		Description:      eventData["description"].(string),
		Price:            eventData["price"].(float64),
		AvailablePortions: int(eventData["available_portions"].(float64)),
		MealCourse:       eventData["meal_course"].(string),
		DietaryCategory:  eventData["dietary_category"].(string),
		AvailableUntil:   availableUntil,
		Location:         geoPoint,
		Address:          address,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		DeletedAt:        deletedAt,
	}

	return dish, nil
}

// Helper functions to parse specific fields
func parseTime(data map[string]interface{}, field string) (time.Time, error) {
	fieldStr, ok := data[field].(string)
	if !ok || fieldStr == "" {
		return time.Time{}, fmt.Errorf("invalid or missing %s in event data", field)
	}
	return time.Parse(time.RFC3339, fieldStr)
}

func parseOptionalTime(data map[string]interface{}, field string) *time.Time {
	fieldStr, ok := data[field].(string)
	if !ok || fieldStr == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, fieldStr)
	if err != nil {
		return nil
	}
	return &t
}


func parseGeoPoint(data map[string]interface{}) (model.Location, error) {
	geoPoint, ok := data["location"].(map[string]interface{})
	if !ok {
		return model.Location{}, fmt.Errorf("invalid or missing location in event data")
	}
	latitude, latOk := geoPoint["latitude"].(float64)
	longitude, lonOk := geoPoint["longitude"].(float64)
	if !latOk || !lonOk {
		return model.Location{}, fmt.Errorf("invalid coordinates in location")
	}
	return model.Location{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}, nil
}

func parseAddress(data map[string]interface{}) (model.Address, error) {
	address, ok := data["address"].(map[string]interface{})
	if !ok {
		return model.Address{}, fmt.Errorf("invalid or missing address in event data")
	}
	return model.Address{
		Street:     address["street"].(string),
		City:       address["city"].(string),
		State:      address["state"].(string),
		PostalCode: address["postal_code"].(string),
		Country:    address["country"].(string),
	}, nil
}
