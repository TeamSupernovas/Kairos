package edu.sjsu.kairos.dishmanagementservice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.Instant;
import java.util.List;
import java.util.UUID;

@Builder
@Data
@AllArgsConstructor
@NoArgsConstructor
public class DishDTO {

    @JsonProperty("dish_id")
    private UUID dishId;

    @JsonProperty("dish_name")
    private String dishName;

    @JsonProperty("chef_id")
    private String chefId;

    @JsonProperty("price")
    private double price;

    @JsonProperty("available_portions")
    private int availablePortions;

    @JsonProperty("description")
    private String description;

    @JsonProperty("meal_course")
    private String mealCourse;

    @JsonProperty("dietary_category")
    private String dietaryCategory;

    @JsonProperty("ingredients")
    private String ingredients;

    @JsonProperty("available_until")
    private Instant availableUntil;

    @JsonProperty("address")
    private AddressDTO address;

    @JsonProperty("location")
    private LocationDTO location;

    @JsonProperty("images")
    private List<ImageDTO> images;

    @JsonProperty("created_at")
    private Instant createdAt;

    @JsonProperty("updated_at")
    private Instant updatedAt;
}
