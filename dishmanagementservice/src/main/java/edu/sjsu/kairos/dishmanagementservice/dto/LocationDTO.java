package edu.sjsu.kairos.dishmanagementservice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Builder
@Data
@AllArgsConstructor
@NoArgsConstructor
public class LocationDTO {

    @JsonProperty("latitude")
    private double latitude;

    @JsonProperty("longitude")
    private double longitude;
}
