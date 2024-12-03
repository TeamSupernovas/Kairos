package mapper

import (
	"dishmanagementservice/internal/dto"
	"dishmanagementservice/internal/model"
)

func ToDishModel(req dto.CreateDishRequest) *model.Dish {
	newDish := model.NewDish(req.DishName, req.ChefID, req.Price, req.AvailablePortions)
	newDish.Description = req.Description
	newDish.MealCourse = req.MealCourse
	newDish.DietaryCategory = req.DietaryCategory
	newDish.AvailableUntil = req.AvailableUntil

	newDish.SetAddress(
		req.Address.Street,
		req.Address.City,
		req.Address.State,
		req.Address.PostalCode,
		req.Address.Country,
	)

	return newDish
}