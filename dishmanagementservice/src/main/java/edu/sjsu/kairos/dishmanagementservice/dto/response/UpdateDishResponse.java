package edu.sjsu.kairos.dishmanagementservice.dto.response;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.UUID;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class UpdateDishResponse {
    @JsonProperty("dish_id")
    private UUID dishId;

    @JsonProperty("message")
    private String message;
}
