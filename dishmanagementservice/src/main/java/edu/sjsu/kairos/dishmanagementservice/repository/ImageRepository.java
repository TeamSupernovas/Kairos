package edu.sjsu.kairos.dishmanagementservice.repository;


import edu.sjsu.kairos.dishmanagementservice.model.Image;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.UUID;

public interface ImageRepository extends JpaRepository<Image, UUID> {

    List<Image> findByDish_DishIdAndDeletedAtIsNull(UUID dishId);

}
