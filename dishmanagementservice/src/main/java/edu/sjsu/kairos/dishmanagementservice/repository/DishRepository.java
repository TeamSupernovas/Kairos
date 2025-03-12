package edu.sjsu.kairos.dishmanagementservice.repository;

import java.util.Optional;
import java.util.UUID;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import edu.sjsu.kairos.dishmanagementservice.model.Dish;

public interface DishRepository extends JpaRepository<Dish, UUID>{

    Optional<Dish> findByDishIdAndDeletedAtIsNull(UUID dishId);
    Page<Dish> findByChefIdAndDeletedAtIsNull(String chefId, Pageable pageable);
}
