package edu.sjsu.kairos.dishmanagementservice.exception;

import lombok.extern.log4j.Log4j2;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.time.LocalDateTime;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.logging.Level;

@RestControllerAdvice
@Log4j2
public class DishManagementServiceExceptionHandler {

    @ExceptionHandler(DishNotFoundException.class)
    public ResponseEntity<Object> handleDishNotFoundException(DishNotFoundException ex) {
        Map<String, Object> errorDetails = new LinkedHashMap<>();
        errorDetails.put("timestamp", LocalDateTime.now());
        errorDetails.put("status", HttpStatus.NOT_FOUND.value());
        errorDetails.put("error", "Dish Not Found");
        errorDetails.put("message", ex.getMessage());
        errorDetails.put("path", "/dishes/{dishId}");
        return new ResponseEntity<>(errorDetails, HttpStatus.NOT_FOUND);
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<Object> handleGlobalException(Exception ex) {
        Map<String, Object> errorDetails = new LinkedHashMap<>();
        errorDetails.put("timestamp", LocalDateTime.now());
        errorDetails.put("status", HttpStatus.INTERNAL_SERVER_ERROR.value());
        errorDetails.put("error", "Internal Server Error");
        errorDetails.put("message", ex.getMessage());
        log.error("exception", ex);
        return new ResponseEntity<>(errorDetails, HttpStatus.INTERNAL_SERVER_ERROR);
    }
}
