CREATE TABLE ratings (
    id SERIAL PRIMARY KEY,               -- Unique ID for each rating
    dish_id VARCHAR(50) NOT NULL,                -- ID of the rated dish
    chef_id VARCHAR(50) NOT NULL,               -- ID of the chef
    user_id VARCHAR(50) NOT NULL,               -- ID of the user giving the rating
    rating INT CHECK (rating BETWEEN 1 AND 5), -- Rating value (1 to 5)
    review_text VARCHAR(250),                    -- Optional review text
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp when rating was created
);
