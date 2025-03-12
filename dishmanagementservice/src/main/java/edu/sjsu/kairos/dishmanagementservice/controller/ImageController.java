package edu.sjsu.kairos.dishmanagementservice.controller;

import edu.sjsu.kairos.dishmanagementservice.dto.ImageDTO;
import edu.sjsu.kairos.dishmanagementservice.service.DishImageService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import java.util.List;
import java.util.UUID;

@CrossOrigin(origins = "*")
@RestController
@RequestMapping("/images")
public class ImageController {

    @Autowired
    private DishImageService imageService;

    @GetMapping("/dish/{dishId}")
    public ResponseEntity<List<ImageDTO>> getImagesByDishId(@PathVariable UUID dishId) {
        List<ImageDTO> images = imageService.getImagesByDishId(dishId);
        return ResponseEntity.ok(images);
    }
}
