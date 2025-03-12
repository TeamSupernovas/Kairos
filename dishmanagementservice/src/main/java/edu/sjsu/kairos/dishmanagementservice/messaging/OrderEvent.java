package edu.sjsu.kairos.dishmanagementservice.messaging;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OrderEvent {
    @JsonProperty("order_id")
    private UUID orderId;

    @JsonProperty("dish_id")
    private UUID dishId;

    @JsonProperty("portions")
    private int portions;
}
