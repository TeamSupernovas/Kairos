package service

import (
	"dishmanagementservice/internal/infrastructure"
	"dishmanagementservice/internal/model"
	"dishmanagementservice/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DishService struct {
	repo *repository.PostgresDishRepository
	geoLocationService *GeoLocationService
	kafkaService *KafkaService
}

func NewDishService(repo *repository.PostgresDishRepository, geoService *GeoLocationService, kafkaService *KafkaService) *DishService {
	return &DishService{
		repo: repo,
		geoLocationService: geoService,
		kafkaService: kafkaService,
	}
}

func (service *DishService) AddDish(dish *model.Dish) error {

	dish.DishID = uuid.NewString()

	if (dish.DishName == "") {
		return fmt.Errorf("dish name is required")
	}

	if (dish.ChefID == "") {
		return fmt.Errorf("chef ID is required")
	}
	if (dish.Price <= 0) {
		return fmt.Errorf("price must be greater than zero")
	}

	if (dish.AvailablePortions <= 0) {
		return fmt.Errorf("available portions must be greater than zero")
	}

	if (!dish.AvailableUntil.IsZero() && dish.AvailableUntil.Before(time.Now())) {
		return fmt.Errorf("available until date/time must be in future")
	}

	address := fmt.Sprintf(
		"%s, %s, %s, %s, %s",
		dish.Address.Street,
		dish.Address.City,
		dish.Address.State,
		dish.Address.PostalCode,
		dish.Address.Country,
	)
	geoPoint, err := service.geoLocationService.CalculateGeoPoint(address, "KairosPlaceIndex")
	if err != nil {
		return fmt.Errorf("failed to derive geolocation: %v", err)
	}
	dish.GeoPoint.Latitude = geoPoint.Latitude
	dish.GeoPoint.Longitude = geoPoint.Longitude

	err = service.repo.CreateDish(dish)
	if err != nil {
		return err
	}

	err = service.kafkaService.PublishEvent(infrastructure.EventDishCreated, dish)
    if err != nil {
	fmt.Printf("Error publishing event: %v\n", err)
	}
	
	return nil
}

func (service *DishService) UpdateDish(dishID string, updatedDish *model.Dish) error {
	if updatedDish.DishName == "" {
		return fmt.Errorf("dish name is required")
	}
	if updatedDish.ChefID == "" {
		return fmt.Errorf("chef ID is required")
	}
	if (updatedDish.Price <= 0) {
		return fmt.Errorf("price must be greater than zero")
	}

	if (updatedDish.AvailablePortions <= 0) {
		return fmt.Errorf("available portions must be greater than zero")
	}

	if (!updatedDish.AvailableUntil.IsZero() && updatedDish.AvailableUntil.Before(time.Now())) {
		return fmt.Errorf("available until date/time must be in future")
	}


	// Call repository to update the dish in the database
	return service.repo.UpdateDish(dishID, updatedDish)
}

func (service *DishService) GetDishByID(dishID string) (*model.Dish, error) {
	if dishID == "" {
		return nil, fmt.Errorf("dish ID is required")
	}
	
	return service.repo.GetDishByID(dishID)
}

func (service *DishService) DeleteDishByID(dishID string) (error) {
	if (dishID == "") {
		return fmt.Errorf("dish ID is required")
	}

	return service.repo.DeleteDishByID(dishID)
}