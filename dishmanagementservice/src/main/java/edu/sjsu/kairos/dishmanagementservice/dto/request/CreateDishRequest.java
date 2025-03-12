package edu.sjsu.kairos.dishmanagementservice.dto.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import edu.sjsu.kairos.dishmanagementservice.dto.AddressDTO;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.Instant;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class CreateDishRequest {

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

    @JsonProperty("address")
    private AddressDTO address;

    @JsonProperty("available_until")
    private Instant availableUntil;
}
