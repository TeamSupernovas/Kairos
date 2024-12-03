package repository

import (
	"database/sql"
	"dishmanagementservice/internal/model"
	"fmt"
	"time"
)

type DishRepository interface {
	CreateDish(dish *model.Dish) (error)
	UpdateDish(dishID string, dish *model.Dish) error
	DeleteDishByID(dishID string) error
	GetDishByID(dishID string) (*model.Dish, error)
}

type PostgresDishRepository struct {
	db *sql.DB
}

func NewPostgresDishRepository(db *sql.DB) *PostgresDishRepository {
	return &PostgresDishRepository{
		db: db,
	}
}

func (repo *PostgresDishRepository) CreateDish(dish *model.Dish) error {
	query := `INSERT INTO dishes (dish_id, dish_name, chef_id, description, price, available_portions, meal_course, dietary_category, available_until, street, city, state, postal_code, country, latitude, longitude, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);`

	_, err := repo.db.Exec(
		query,
		dish.DishID,
		dish.DishName,
		dish.ChefID,
		dish.Description,
		dish.Price,
		dish.AvailablePortions,
		dish.MealCourse,
		dish.DietaryCategory,
		dish.AvailableUntil,
		dish.Address.Street,
		dish.Address.City,
		dish.Address.State,
		dish.Address.PostalCode,
		dish.Address.Country,
		dish.GeoPoint.Latitude,
		dish.GeoPoint.Longitude,
		dish.CreatedAt,
		dish.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresDishRepository) UpdateDish(dishID string, dish *model.Dish) error {
	query := `
		UPDATE dishes 
		SET dish_name = $1, 
			chef_id = $2, 
			description = $3, 
			price = $4, 
			available_portions = $5, 
			meal_course = $6, 
			dietary_category = $7, 
			available_until = $8, 
			street = $9, 
			city = $10, 
			state = $11, 
			postal_code = $12, 
			country = $13, 
			latitude = $14, 
			longitude = $15, 
			updated_at = $16
		WHERE dish_id = $17;
	`

	_, err := repo.db.Exec(
		query,
		dish.DishName,
		dish.ChefID,
		dish.Description,
		dish.Price,
		dish.AvailablePortions,
		dish.MealCourse,
		dish.DietaryCategory,
		dish.AvailableUntil,
		dish.Address.Street,
		dish.Address.City,
		dish.Address.State,
		dish.Address.PostalCode,
		dish.Address.Country,
		dish.GeoPoint.Latitude,
		dish.GeoPoint.Longitude,
		time.Now(),
		dishID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresDishRepository) GetDishByID(dishID string) (*model.Dish, error) {
	query := `
		SELECT dish_id, dish_name, chef_id, description, price, available_portions, meal_course, dietary_category, available_until, street, city, state, postal_code, country, latitude, longitude, created_at, updated_at
		FROM dishes
		WHERE dish_id = $1
		AND deleted_at is NULL;
	`

	var dish model.Dish
	err := repo.db.QueryRow(query, dishID).Scan(
		&dish.DishID,
		&dish.DishName,
		&dish.ChefID,
		&dish.Description,
		&dish.Price,
		&dish.AvailablePortions,
		&dish.MealCourse,
		&dish.DietaryCategory,
		&dish.AvailableUntil,
		&dish.Address.Street,
		&dish.Address.City,
		&dish.Address.State,
		&dish.Address.PostalCode,
		&dish.Address.Country,
		&dish.GeoPoint.Latitude,
		&dish.GeoPoint.Longitude,
		&dish.CreatedAt,
		&dish.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no dish found with the provided ID")
	} else if err != nil {
		return nil, err
	}
	return &dish, nil
}

func (repo *PostgresDishRepository) DeleteDishByID(dishID string) error {
	query := `UPDATE dishes SET deleted_at = NOW() WHERE dish_id = $1`

	// Execute the delete query
	_, err := repo.db.Exec(query, dishID)
	if (err != nil) {
		return fmt.Errorf("could not soft delete dish: %v", err)
	}
	return nil
}
