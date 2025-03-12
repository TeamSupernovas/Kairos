package edu.sjsu.kairos.dishmanagementservice.exception;

public class DishNotFoundException extends RuntimeException {
    public DishNotFoundException(String message) {
        super(message);
    }
}