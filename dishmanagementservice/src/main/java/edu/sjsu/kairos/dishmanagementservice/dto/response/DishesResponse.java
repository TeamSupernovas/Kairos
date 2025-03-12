package edu.sjsu.kairos.dishmanagementservice.dto.response;

import com.fasterxml.jackson.annotation.JsonProperty;
import edu.sjsu.kairos.dishmanagementservice.dto.DishDTO;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.util.List;

@Builder
@Data
@AllArgsConstructor
@NoArgsConstructor
public class DishesResponse {
    @JsonProperty("dishes")
    private List<DishDTO> dishes;
    private int currentPage;
    private long totalItems;
    private int totalPages;
    private int pageSize;
}
