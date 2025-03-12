package edu.sjsu.kairos.dishmanagementservice.messaging;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;
import edu.sjsu.kairos.dishmanagementservice.util.DishEventType;
import edu.sjsu.kairos.dishmanagementservice.util.DishReservationStatus;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.Instant;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@JsonSerialize
@JsonInclude(JsonInclude.Include.NON_NULL)
public class DishReservationEvent {
    @JsonProperty("dish_id")
    private UUID dishId;

    @JsonProperty("event_type")
    private DishEventType eventType;

    @JsonProperty("event_time")
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss'Z'", timezone = "UTC")
    private Instant eventTime;

    @JsonProperty("order_id")
    private UUID orderId;

    @JsonProperty("chef_id")
    private String chefId;

    @JsonProperty("requested_portions")
    private int requestedPortions;

    @JsonProperty("confirmed_portions")
    private int confirmedPortions;

    @JsonProperty("remaining_portions")
    private int remainingPortions;

    @JsonProperty("status")
    private DishReservationStatus status;
}
