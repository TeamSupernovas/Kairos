package edu.sjsu.kairos.dishmanagementservice.mapper;

import edu.sjsu.kairos.dishmanagementservice.dto.AddressDTO;
import edu.sjsu.kairos.dishmanagementservice.dto.DishDTO;
import edu.sjsu.kairos.dishmanagementservice.dto.LocationDTO;
import edu.sjsu.kairos.dishmanagementservice.dto.request.CreateDishRequest;
import edu.sjsu.kairos.dishmanagementservice.dto.request.UpdateDishRequest;
import edu.sjsu.kairos.dishmanagementservice.dto.response.CreateDishResponse;
import edu.sjsu.kairos.dishmanagementservice.dto.response.DishResponse;
import edu.sjsu.kairos.dishmanagementservice.dto.response.UpdateDishResponse;
import edu.sjsu.kairos.dishmanagementservice.model.Address;
import edu.sjsu.kairos.dishmanagementservice.model.Dish;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
public class DishMapper {

    @Autowired
    private AddressMapper addressMapper;

    @Autowired
    private LocationMapper locationMapper;

    public DishDTO toDishDTO(Dish dish) {
        AddressDTO addressDTO = null;
        if (dish.getAddress() != null) {
            addressDTO = addressMapper.toAddressDTO(dish.getAddress());
        }

        LocationDTO locationDTO = null;
        if (dish.getLocation() != null) {
            locationDTO = locationMapper.toLocationDTO(dish.getLocation());
        }

        return DishDTO.builder()
                .dishId(dish.getDishId())
                .dishName(dish.getDishName())
                .chefId(dish.getChefId())
                .price(dish.getPrice())
                .availablePortions(dish.getAvailablePortions())
                .description(dish.getDescription())
                .mealCourse(dish.getMealCourse())
                .dietaryCategory(dish.getDietaryCategory())
                .ingredients(dish.getIngredients())
                .availableUntil(dish.getAvailableUntil())
                .address(addressDTO)
                .location(locationDTO)
                .createdAt(dish.getCreatedAt())
                .updatedAt(dish.getUpdatedAt())
                .build();
    }

    public Dish toDish(CreateDishRequest createDishRequest){
        Address address = null;
        if (createDishRequest.getAddress() != null) {
            address = addressMapper.toAddress(createDishRequest.getAddress());
        }
         addressMapper.toAddress(createDishRequest.getAddress());
        return Dish.builder()
                .dishName(createDishRequest.getDishName())
                .chefId(createDishRequest.getChefId())
                .price(createDishRequest.getPrice())
                .availablePortions(createDishRequest.getAvailablePortions())
                .description(createDishRequest.getDescription())
                .mealCourse(createDishRequest.getMealCourse())
                .dietaryCategory(createDishRequest.getDietaryCategory())
                .ingredients(createDishRequest.getIngredients())
                .availableUntil(createDishRequest.getAvailableUntil())
                .address(address)
                .build();
    }

    public DishResponse toDishResponse(Dish dish) {
        DishDTO dishDTO = toDishDTO(dish);
        return DishResponse.builder()
                .dish(dishDTO)
                .build();
    }

    public CreateDishResponse toCreateDishResponse(Dish dish) {
        return CreateDishResponse.builder()
                .dishId(dish.getDishId())
                .message("Dish Created Successfully!")
                .build();
    }

    public UpdateDishResponse toUpdateDishResponse(Dish dish) {
        return UpdateDishResponse.builder()
                .dishId(dish.getDishId())
                .message("Dish Updated Successfully!")
                .build();
    }

    public void mapUpdateDishRequestToDish(UpdateDishRequest updateDishRequest, Dish dish) {
        if (updateDishRequest != null && dish != null) {
            dish.setDishName(updateDishRequest.getDishName());
            dish.setChefId(updateDishRequest.getChefId());
            dish.setPrice(updateDishRequest.getPrice());
            dish.setAvailablePortions(updateDishRequest.getAvailablePortions());
            dish.setDescription(updateDishRequest.getDescription());
            dish.setMealCourse(updateDishRequest.getMealCourse());
            dish.setDietaryCategory(updateDishRequest.getDietaryCategory());
            dish.setIngredients(updateDishRequest.getIngredients());
            dish.setAvailableUntil(updateDishRequest.getAvailableUntil());

            if (updateDishRequest.getAddress() != null) {
                AddressDTO addressDTO = updateDishRequest.getAddress();
                Address address = dish.getAddress();
                address.setStreet(addressDTO.getStreet());
                address.setCity(addressDTO.getCity());
                address.setPostalCode(addressDTO.getPostalCode());
                address.setState(addressDTO.getState());
                address.setCountry(addressDTO.getCountry());
            }
        }
    }

    public List<DishDTO> toDishDTOList(List<Dish> dishes) {
        return dishes.stream().map(this::toDishDTO).toList();
    }
}
