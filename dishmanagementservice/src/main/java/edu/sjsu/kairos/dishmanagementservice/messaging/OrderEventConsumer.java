package edu.sjsu.kairos.dishmanagementservice.messaging;

import edu.sjsu.kairos.dishmanagementservice.mapper.DishMapper;
import edu.sjsu.kairos.dishmanagementservice.model.Dish;
import edu.sjsu.kairos.dishmanagementservice.repository.DishRepository;
import edu.sjsu.kairos.dishmanagementservice.util.DishEventType;
import edu.sjsu.kairos.dishmanagementservice.util.DishReservationStatus;
import jakarta.transaction.Transactional;
import lombok.extern.log4j.Log4j2;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.util.Optional;
import java.util.UUID;

@Service
@Log4j2
public class OrderEventConsumer {

    @Autowired
    private DishRepository dishRepository;

    @Autowired
    private DishEventProducer dishEventProducer;

    @Autowired
    private DishMapper dishMapper;

    @KafkaListener(topics = "${kairos.orderservice.kafka.topicname.order.placed}", containerFactory = "orderEventKafkaListenerContainerFactory")
    public void consumeOrderPlacedEvent(OrderEvent orderEvent) {
        System.out.println("Received OrderEvent: OrderID=" + orderEvent.getOrderId());
        log.info("Received OrderEvent: OrderID={}, DishID={}, Portions={}",
                orderEvent.getOrderId(), orderEvent.getDishId(), orderEvent.getPortions());
        processOrder(orderEvent);
    }

    private void processOrder(OrderEvent orderEvent) {
        // Business logic to check dish availability and update portions
        log.info("Processing order for DishID={}, Requested Portions={}", orderEvent.getDishId(), orderEvent.getPortions());

        boolean isDishAvailable = checkDishAvailability(orderEvent.getDishId(), orderEvent.getPortions());

        if (isDishAvailable) {
            log.info("DishID={} is available. Reserving portions...", orderEvent.getDishId());
            reserveDish(orderEvent);
        } else {
            log.warn("DishID={} is NOT available for {} portions. Sending rejection event.", orderEvent.getDishId(), orderEvent.getPortions());
            DishReservationEvent dishReservationEvent = DishReservationEvent.builder()
                    .dishId(orderEvent.getDishId())
                    .orderId(orderEvent.getOrderId())
                    .requestedPortions(orderEvent.getPortions())
                    .eventTime(Instant.now())
                    .eventType(DishEventType.RESERVATION_STATUS)
                    .status(DishReservationStatus.REJECTED)
                    .build();
            dishEventProducer.sendDishReservationEvent(dishReservationEvent);
        }
    }

    @Transactional
    public void reserveDish(OrderEvent orderEvent) {
        UUID dishId = orderEvent.getDishId();
        int requestedPortions = orderEvent.getPortions();

        Optional<Dish> dishOptional = dishRepository.findByDishIdAndDeletedAtIsNull(dishId);

        if (dishOptional.isEmpty()) {
            log.warn("DishID={} not found or deleted. Unable to reserve portions.", dishId);
            DishReservationEvent dishReservationEvent = DishReservationEvent.builder()
                    .eventType(DishEventType.RESERVATION_STATUS)
                    .eventTime(Instant.now())
                    .confirmedPortions(0)
                    .remainingPortions(0)
                    .requestedPortions(requestedPortions)
                    .orderId(orderEvent.getOrderId())
                    .dishId(dishId)
                    .status(DishReservationStatus.REJECTED)
                    .build();
            dishEventProducer.sendDishReservationEvent(dishReservationEvent);
            return;
        }

        Dish dish = dishOptional.get();

        if (dish.getAvailablePortions() < requestedPortions) {
            log.warn("DishID={} does not have enough portions. Requested={}, Available={}",
                    dishId, requestedPortions, dish.getAvailablePortions());
            DishReservationEvent dishReservationEvent = DishReservationEvent.builder()
                    .dishId(dishId)
                    .orderId(orderEvent.getOrderId())
                    .chefId(dish.getChefId())
                    .requestedPortions(requestedPortions)
                    .confirmedPortions(0)
                    .remainingPortions(dish.getAvailablePortions())
                    .eventTime(Instant.now())
                    .eventType(DishEventType.RESERVATION_STATUS)
                    .status(DishReservationStatus.REJECTED)
                    .build();
            dishEventProducer.sendDishReservationEvent(dishReservationEvent);
            return;
        }

        // Deduct portions
        dish.setAvailablePortions(dish.getAvailablePortions() - requestedPortions);
        Dish savedDish = dishRepository.save(dish); // Save updated dish

        log.info("Successfully reserved {} portions for DishID={}. Remaining portions={}",
                requestedPortions, dishId, dish.getAvailablePortions());

        DishReservationEvent dishReservationEvent = DishReservationEvent.builder()
                .dishId(dishId)
                .orderId(orderEvent.getOrderId())
                .chefId(savedDish.getChefId())
                .requestedPortions(requestedPortions)
                .confirmedPortions(requestedPortions)
                .remainingPortions(savedDish.getAvailablePortions())
                .eventTime(Instant.now())
                .eventType(DishEventType.RESERVATION_STATUS)
                .status(DishReservationStatus.CONFIRMED)
                .build();
        dishEventProducer.sendDishReservationEvent(dishReservationEvent);


        DishEvent dishEvent = updateDishEvent(savedDish, DishEventType.UPDATED);
        dishEventProducer.sendDishEvent(dishEvent);
    }

    public boolean checkDishAvailability(UUID dishId, int requestedPortions) {
        Optional<Dish> dishOptional = dishRepository.findByDishIdAndDeletedAtIsNull(dishId);

        if (dishOptional.isEmpty()) {
            log.warn("DishID={} not found or deleted.", dishId);
            return false;
        }

        Dish dish = dishOptional.get();

        if (dish.getAvailablePortions() >= requestedPortions) {
            log.info("DishID={} has enough portions available. Requested={}, Available={}",
                    dishId, requestedPortions, dish.getAvailablePortions());
            return true; // Portions are available
        }

        log.warn("DishID={} does not have enough portions. Requested={}, Available={}",
                dishId, requestedPortions, dish.getAvailablePortions());
        return false;
    }

    private DishEvent updateDishEvent(Dish savedDish, DishEventType dishEventType) {
        Address address = Address.builder()
                .street(savedDish.getAddress().getStreet())
                .city(savedDish.getAddress().getCity())
                .state(savedDish.getAddress().getState())
                .postalCode(savedDish.getAddress().getPostalCode())
                .country(savedDish.getAddress().getCountry())
                .build();

        Location location = Location.builder()
                .latitude(savedDish.getLocation().getLatitude())
                .longitude(savedDish.getLocation().getLongitude())
                .build();

        DishEvent dishEvent = DishEvent.builder()
                .dishId(savedDish.getDishId())
                .eventType(dishEventType)
                .eventTime(Instant.now())
                .dishName(savedDish.getDishName())
                .chefId(savedDish.getChefId())
                .price(savedDish.getPrice())
                .availablePortions(savedDish.getAvailablePortions())
                .description(savedDish.getDescription())
                .mealCourse(savedDish.getMealCourse())
                .dietaryCategory(savedDish.getDietaryCategory())
                .ingredients(savedDish.getIngredients())
                .availableUntil(savedDish.getAvailableUntil())
                .address(address)
                .location(location)
                .createdAt(savedDish.getCreatedAt())
                .updatedAt(savedDish.getUpdatedAt())
                .deletedAt(savedDish.getDeletedAt())
                .build();
        return dishEvent;
    }

}
