package edu.sjsu.kairos.dishmanagementservice.messaging;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Builder
@Data
@AllArgsConstructor
@NoArgsConstructor
@JsonSerialize
@JsonInclude(JsonInclude.Include.NON_NULL)
public class Location {

    @JsonProperty("latitude")
    private double latitude;

    @JsonProperty("longitude")
    private double longitude;
}
