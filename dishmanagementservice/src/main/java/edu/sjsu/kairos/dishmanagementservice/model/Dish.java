package edu.sjsu.kairos.dishmanagementservice.model;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.Instant;
import java.util.List;
import java.util.UUID;


@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
@Entity
@Table(name = "dishes")
public class Dish {
	
	@Id
	@GeneratedValue(strategy = GenerationType.UUID)
    @Column(name="dish_id")
    private UUID dishId;
	
	@Column(name="dish_name", nullable = false, length = 255)
    private String dishName;
	
	@Column(name="chef_id", nullable = false, length = 50)
    private String chefId; 
	
	@Column(name="price", nullable = false)
    private double price;
	
	@Column(name="available_portions", nullable = false)
    private int availablePortions;
	
	@Column(name="description", columnDefinition = "TEXT")
    private String description;
	
	@Column(name="meal_course", length = 50)
    private String mealCourse;
	
	@Column(name="dietary_category", length = 50)
    private String dietaryCategory;
	
	@Column(name="ingredients", columnDefinition = "TEXT")
    private String ingredients;
	
	@Column(name="available_until")
    private Instant availableUntil;

    @Embedded
    private Address address;

    @Embedded
    private Location location;

    @OneToMany(mappedBy = "dish", cascade = CascadeType.ALL)
    private List<Image> images;

    @Column(name="created_at", nullable = false, updatable = false)
    private Instant createdAt;

    @Column(name="updated_at", nullable = false)
    private Instant updatedAt;

    @Column(name="deleted_at")
    private Instant deletedAt;

    @PrePersist
    protected void onCreate() {
    	this.createdAt = Instant.now();
        this.updatedAt = Instant.now();
    }

    @PreUpdate
    protected void onUpdate() {
    	this.updatedAt = Instant.now();
    }

    public void markAsDeleted() {
        this.deletedAt = Instant.now();
    }
    
    public boolean isDeleted() {
        return this.deletedAt != null;
    }
}
