CREATE TABLE orders (
    order_id  VARCHAR(50) PRIMARY KEY NOT NULL,
    user_id  VARCHAR(50) NOT NULL,
    chef_id  VARCHAR(50) NOT NULL,
    total_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    pickup_time TIMESTAMPTZ,                            -- Nullable
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    canceled_at TIMESTAMPTZ,                            -- Nullable
    completed_at TIMESTAMPTZ,                           -- Nullable
    order_status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (
        order_status IN ('pending', 'confirmed', 'ready','canceled', 'completed')
    ),
    deleted_at TIMESTAMPTZ                              -- Nullable
);

CREATE TABLE order_items (
    order_item_id  VARCHAR(50) PRIMARY KEY NOT NULL,
    order_id  VARCHAR(50) NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    dish_id  VARCHAR(50) NOT NULL,
    dish_order_status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (
        dish_order_status IN ('pending', 'confirmed', 'ready', 'canceled', 'completed')
    ),
    quantity INT NOT NULL,
    price_per_unit DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ                              -- Nullable
);