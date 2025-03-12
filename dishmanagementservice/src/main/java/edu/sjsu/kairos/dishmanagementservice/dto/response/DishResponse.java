package edu.sjsu.kairos.dishmanagementservice.dto.response;

import com.fasterxml.jackson.annotation.JsonProperty;
import edu.sjsu.kairos.dishmanagementservice.dto.DishDTO;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Builder
@Data
@AllArgsConstructor
@NoArgsConstructor
public class DishResponse {
    @JsonProperty("dish")
    private DishDTO dish;

}
