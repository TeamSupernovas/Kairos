package edu.sjsu.kairos.dishmanagementservice.messaging;

import edu.sjsu.kairos.dishmanagementservice.util.DishEventType;
import lombok.extern.log4j.Log4j2;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Service;
import java.util.concurrent.CompletableFuture;


@Service
@Log4j2
public class DishEventProducer {

    private final KafkaTemplate<String, Object> kafkaTemplate;

    @Value("${kairos.dishmanagementservice.kafka.topicname.dish.reservationstatus}")
    private String reservationStatusEventTopicName;

    @Value("${kairos.dishmanagementservice.kafka.topicname.dish.created}")
    private String dishCreatedEventTopicName;

    @Value("${kairos.dishmanagementservice.kafka.topicname.dish.updated}")
    private String dishEventUpdatedTopicName;

    @Value("${kairos.dishmanagementservice.kafka.topicname.dish.deleted}")
    private String dishDeletedEventTopicName;

    @Autowired
    public DishEventProducer(KafkaTemplate<String, Object> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    public void sendDishEvent(DishEvent dishEvent) {
        String topicName = getTopicForEventType(dishEvent.getEventType());
        String partitionKey = dishEvent.getDishId().toString();

        CompletableFuture<SendResult<String, Object>> future = kafkaTemplate.send(topicName, partitionKey, dishEvent);
        log.info("Sending DishEvent={} to topic={}: DishID={}", dishEvent.getEventType(), topicName, dishEvent.getDishId());
        future.whenComplete((result, ex) -> {
            if (ex == null) {
                log.info("SUCCESS: DishEvent={} sent: DishID={}: Offset={}", dishEvent.getEventType(), dishEvent.getDishId(),
                        result.getRecordMetadata().offset());
            } else {
                log.error("FAILURE: Unable to send DishEvent={}: DishID={}: Error={}", dishEvent.getEventType(), dishEvent.getDishId(),
                        ex.getMessage());
            }
        });
    }

    public void sendDishReservationEvent(DishReservationEvent dishReservationEvent) {
        String partitionKey = dishReservationEvent.getDishId().toString();
        CompletableFuture<SendResult<String, Object>> future = kafkaTemplate.send(reservationStatusEventTopicName, partitionKey, dishReservationEvent);
        log.info("Sending DishReservationEvent={} to topic={} DishID={}, Status={}",
                dishReservationEvent.getEventType(), reservationStatusEventTopicName, dishReservationEvent.getDishId(), dishReservationEvent.getStatus());

        future.whenComplete((result, ex) -> {
            if (ex == null) {
                log.info("SUCCESS: DishReservationEvent={} sent: DishID={}, Offset={}",
                        dishReservationEvent.getEventType(), dishReservationEvent.getDishId(), result.getRecordMetadata().offset());
            } else {
                log.error("FAILURE: Unable to send DishReservationEvent={}: DishID={}, Error={}",
                        dishReservationEvent.getEventType(), dishReservationEvent.getDishId(), ex.getMessage());
            }
        });
    }


    private String getTopicForEventType(DishEventType eventType) {
        return switch (eventType) {
            case CREATED -> dishCreatedEventTopicName;
            case UPDATED -> dishEventUpdatedTopicName;
            case DELETED -> dishDeletedEventTopicName;
            default -> throw new IllegalArgumentException("Unknown DishEventType: " + eventType);
        };
    }
}

