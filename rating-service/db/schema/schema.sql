CREATE TABLE ratings (
    id SERIAL PRIMARY KEY,                                -- Unique ID for each rating
    dish_id VARCHAR(50) NOT NULL,                         -- ID of the rated dish
    dish_name VARCHAR(100) NOT NULL,                      -- Name of the dish (denormalized)
    chef_id VARCHAR(50) NOT NULL,                         -- ID of the chef
    chef_name VARCHAR(100) NOT NULL,                      -- Name of the chef (denormalized)
    user_id VARCHAR(50) NOT NULL,                         -- ID of the user giving the rating
    user_name VARCHAR(100) NOT NULL,                      -- Name of the user (denormalized)
    rating INT CHECK (rating BETWEEN 1 AND 5),            -- Rating value (1 to 5)
    review_text VARCHAR(250),                             -- Optional review text
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP        -- Timestamp when rating was created
);
