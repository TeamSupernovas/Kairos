package edu.sjsu.kairos.dishmanagementservice.messaging;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;
import edu.sjsu.kairos.dishmanagementservice.util.DishEventType;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.Instant;
import java.util.List;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@JsonSerialize
@JsonInclude(JsonInclude.Include.NON_NULL)
public class DishEvent {
    @JsonProperty("dish_id")
    private UUID dishId;

    @JsonProperty("event_type")
    private DishEventType eventType;

    @JsonProperty("event_time")
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss'Z'", timezone = "UTC")
    private Instant eventTime;

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
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss'Z'", timezone = "UTC")
    private Instant availableUntil;

    @JsonProperty("address")
    private Address address;

    @JsonProperty("location")
    private Location location;

    @JsonProperty("created_at")
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss'Z'", timezone = "UTC")
    private Instant createdAt;

    @JsonProperty("updated_at")
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss'Z'", timezone = "UTC")
    private Instant updatedAt;

    @JsonProperty("deleted_at")
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss'Z'", timezone = "UTC")
    private Instant deletedAt;
}
