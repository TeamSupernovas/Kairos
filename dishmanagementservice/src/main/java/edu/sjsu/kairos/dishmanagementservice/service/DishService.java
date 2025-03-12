package edu.sjsu.kairos.dishmanagementservice.service;

import edu.sjsu.kairos.dishmanagementservice.dto.DishDTO;
import edu.sjsu.kairos.dishmanagementservice.dto.ImageDTO;
import edu.sjsu.kairos.dishmanagementservice.dto.request.CreateDishRequest;
import edu.sjsu.kairos.dishmanagementservice.dto.request.UpdateDishRequest;
import edu.sjsu.kairos.dishmanagementservice.dto.response.CreateDishResponse;
import edu.sjsu.kairos.dishmanagementservice.dto.response.DishResponse;
import edu.sjsu.kairos.dishmanagementservice.dto.response.DishesResponse;
import edu.sjsu.kairos.dishmanagementservice.dto.response.UpdateDishResponse;
import edu.sjsu.kairos.dishmanagementservice.exception.DishNotFoundException;
import edu.sjsu.kairos.dishmanagementservice.mapper.DishMapper;
import edu.sjsu.kairos.dishmanagementservice.messaging.DishEvent;
import edu.sjsu.kairos.dishmanagementservice.messaging.DishEventProducer;
import edu.sjsu.kairos.dishmanagementservice.model.Address;
import edu.sjsu.kairos.dishmanagementservice.model.Image;
import edu.sjsu.kairos.dishmanagementservice.model.Dish;
import edu.sjsu.kairos.dishmanagementservice.model.Location;
import edu.sjsu.kairos.dishmanagementservice.repository.DishRepository;
import edu.sjsu.kairos.dishmanagementservice.repository.ImageRepository;
import edu.sjsu.kairos.dishmanagementservice.util.DishEventType;
import jakarta.transaction.Transactional;
import lombok.extern.log4j.Log4j2;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.time.Instant;
import java.util.Collections;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
@Log4j2
public class DishService {

    @Autowired
    private DishRepository dishRepository;

    @Autowired
    private DishMapper dishMapper;

    @Autowired
    private GeoLocationService geoLocationService;

    @Autowired
    private DishEventProducer dishEventProducer;

    @Autowired
    private DishImageService imageService;

    @Autowired
    private ImageRepository imageRepository;

    @Transactional
    public CreateDishResponse createDish(CreateDishRequest createDishRequest) {
        Dish dish = dishMapper.toDish(createDishRequest);

        String fullAddress = getAddressString(dish.getAddress());

        double[] coordinates = geoLocationService.getLatLongFromAddress(fullAddress);
        Location location = Location.builder()
                .latitude(coordinates[0])
                .longitude(coordinates[1])
                .build();
        dish.setLocation(location);

        Dish savedDish = dishRepository.save(dish);

        DishEvent dishEvent = createDishEvent(savedDish, DishEventType.CREATED);
        dishEventProducer.sendDishEvent(dishEvent);

        return dishMapper.toCreateDishResponse(savedDish);
    }

    @Transactional
    public CreateDishResponse createDishWithAttachment(CreateDishRequest createDishRequest, MultipartFile file) throws IOException {
        Dish dish = dishMapper.toDish(createDishRequest);
        if (file != null) {
            Image createdImage = imageService.createImage(file);
            createdImage.setDish(dish);
            dish.setImages(Collections.singletonList(createdImage));
        }
        String fullAddress = getAddressString(dish.getAddress());

        double[] coordinates = geoLocationService.getLatLongFromAddress(fullAddress);
        Location location = Location.builder()
                .latitude(coordinates[0])
                .longitude(coordinates[1])
                .build();
        dish.setLocation(location);

        Dish savedDish = dishRepository.save(dish);

        DishEvent dishEvent = createDishEvent(savedDish, DishEventType.CREATED);
        dishEventProducer.sendDishEvent(dishEvent);

        return dishMapper.toCreateDishResponse(savedDish);
    }

    @Transactional
    public UpdateDishResponse updateDish(UUID dishId, UpdateDishRequest updateDishRequest) {
        Dish existingDish = dishRepository.findByDishIdAndDeletedAtIsNull(dishId)
                .orElseThrow(() -> new DishNotFoundException("Dish not found with ID: " + dishId));

        dishMapper.mapUpdateDishRequestToDish(updateDishRequest, existingDish);
        String fullAddress = getAddressString(existingDish.getAddress());
        double[] coordinates = geoLocationService.getLatLongFromAddress(fullAddress);
        existingDish.getLocation().setLatitude(coordinates[0]);
        existingDish.getLocation().setLongitude(coordinates[1]);

        // Save updated dish
        Dish updatedDish = dishRepository.save(existingDish);

        // Publish Kafka event
        DishEvent dishEvent = createDishEvent(updatedDish, DishEventType.UPDATED);
        dishEventProducer.sendDishEvent(dishEvent);

        return dishMapper.toUpdateDishResponse(updatedDish);
    }


    public DishResponse getDishById(UUID dishId) {
        Dish dish = dishRepository.findByDishIdAndDeletedAtIsNull(dishId)
                .orElseThrow(() -> new DishNotFoundException("Dish not found with ID: " + dishId));
        DishResponse dishResponse = dishMapper.toDishResponse(dish);
        List<ImageDTO> images = imageService.getImagesByDishId(dish.getDishId());
       // log.info("Image URLS : {} ", images.get(0));
        dishResponse.getDish().setImages(images);
        return dishResponse;
    }

    public DishesResponse getActiveDishesByChefId(String chefId, int page, int pageSize) {
        Pageable pageable = PageRequest.of(page, pageSize);
        Page<Dish> dishesPage = dishRepository.findByChefIdAndDeletedAtIsNull(chefId, pageable);

        List<DishDTO> dishDTOS = dishesPage.getContent().stream()
                .map(dish -> {
                    DishDTO dishDTO = dishMapper.toDishDTO(dish);
                    List<ImageDTO> images = imageService.getImagesByDishId(dish.getDishId());
                    dishDTO.setImages(images);
                    return dishDTO;
                })
                .collect(Collectors.toList());

        DishesResponse dishesResponse = DishesResponse.builder()
                .dishes(dishDTOS)
                .currentPage(dishesPage.getNumber())
                .pageSize(dishesPage.getSize())
                .totalItems(dishesPage.getTotalElements())
                .totalPages(dishesPage.getTotalPages())
                .build();


        return dishesResponse;
    }

    @Transactional
    public void deleteDish(UUID dishId) {
        Dish existingDish = dishRepository.findById(dishId)
                .orElseThrow(() -> new DishNotFoundException("Dish not found with ID: " + dishId));

        existingDish.markAsDeleted();
        dishRepository.save(existingDish);

        existingDish.getImages().forEach(Image::markAsDeleted);
        imageRepository.saveAll(existingDish.getImages());


        DishEvent dishEvent = createDishEvent(existingDish, DishEventType.DELETED);
        dishEventProducer.sendDishEvent(dishEvent);

    }

    private DishEvent createDishEvent(Dish savedDish, DishEventType dishEventType) {
        edu.sjsu.kairos.dishmanagementservice.messaging.Address address = edu.sjsu.kairos.dishmanagementservice.messaging.Address.builder()
                .street(savedDish.getAddress().getStreet())
                .city(savedDish.getAddress().getCity())
                .state(savedDish.getAddress().getState())
                .postalCode(savedDish.getAddress().getPostalCode())
                .country(savedDish.getAddress().getCountry())
                .build();

        edu.sjsu.kairos.dishmanagementservice.messaging.Location location = edu.sjsu.kairos.dishmanagementservice.messaging.Location.builder()
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

    private static String getAddressString(Address address) {
        return String.join(", ",
                address.getStreet(),
                address.getCity(),
                address.getState(),
                address.getPostalCode(),
                address.getCountry()
        );
    }
}
