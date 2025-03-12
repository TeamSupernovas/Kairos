package edu.sjsu.kairos.dishmanagementservice.controller;

import com.fasterxml.jackson.annotation.JsonProperty;
import edu.sjsu.kairos.dishmanagementservice.dto.request.CreateDishRequest;
import edu.sjsu.kairos.dishmanagementservice.dto.request.UpdateDishRequest;
import edu.sjsu.kairos.dishmanagementservice.dto.response.CreateDishResponse;
import edu.sjsu.kairos.dishmanagementservice.dto.response.DishResponse;
import edu.sjsu.kairos.dishmanagementservice.dto.response.DishesResponse;
import edu.sjsu.kairos.dishmanagementservice.dto.response.UpdateDishResponse;
import edu.sjsu.kairos.dishmanagementservice.service.DishService;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.UUID;

@CrossOrigin(origins = "*")
@RestController
@RequestMapping("/dishes")
public class DishController {

    @Autowired
    DishService dishService;

    /*@PostMapping
    public ResponseEntity<CreateDishResponse> createDish(@RequestBody CreateDishRequest createDishRequest) {
        CreateDishResponse response = dishService.createDish(createDishRequest);
        return ResponseEntity.ok(response);
    }*/

    @PostMapping(consumes = {MediaType.MULTIPART_FORM_DATA_VALUE})
    public ResponseEntity<CreateDishResponse> createDish(@Valid @RequestPart(value="dish", required=true) CreateDishRequest createDishRequest,
                                                               @RequestPart(value="image", required=false) MultipartFile file) throws IOException {
        CreateDishResponse response = dishService.createDishWithAttachment(createDishRequest, file);
        return ResponseEntity.ok(response);
    }

    // Update Dish
    @PutMapping("/{dishId}")
    public ResponseEntity<UpdateDishResponse> updateDish(@PathVariable UUID dishId, @RequestBody UpdateDishRequest request) {
        UpdateDishResponse response = dishService.updateDish(dishId, request);
        return ResponseEntity.ok(response);
    }

    // Get Dish by ID
    @GetMapping("/{dishId}")
    public ResponseEntity<DishResponse> getDishByDishId(@PathVariable UUID dishId) {
        DishResponse response = dishService.getDishById(dishId);
        return ResponseEntity.ok(response);
    }

    @GetMapping("/chef/{chefId}")
    public ResponseEntity<DishesResponse> getDishesByChef(
            @PathVariable String chefId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int pageSize) {

        DishesResponse response = dishService.getActiveDishesByChefId(chefId, page, pageSize);
        return ResponseEntity.ok(response);
    }

    // Delete Dish
    @PatchMapping("/{dishId}")
    public ResponseEntity<Void> deleteDish(@PathVariable UUID dishId) {
        dishService.deleteDish(dishId);
        return ResponseEntity.noContent().build();
    }
}
