CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    userId INT NOT NULL,                               -- Unique identifier for each order
    store_id INT NOT NULL,                             -- Store ID, required field
    merchant_order_id TEXT,                            -- Optional field for merchant's order ID
    recipient_name TEXT NOT NULL,                      -- Recipient's name, required field
    recipient_phone TEXT NOT NULL,                     -- Recipient's phone number, required field
    recipient_address TEXT NOT NULL,                   -- Recipient's address, required field
    recipient_city INT NOT NULL,                       -- Recipient's city ID, required field
    recipient_zone INT NOT NULL,                       -- Recipient's zone ID, required field
    recipient_area INT NOT NULL,                       -- Recipient's area ID, required field
    delivery_type INT NOT NULL,                        -- Delivery type ID, required field
    item_type INT NOT NULL,                            -- Item type ID, required field
    special_instruction TEXT,                          -- Optional field for special instructions
    item_quantity INT NOT NULL,                        -- Item quantity, required field
    item_weight FLOAT NOT NULL,                        -- Item weight in kilograms, required field
    amount_to_collect FLOAT NOT NULL,                  -- Amount to collect, required field
    item_description TEXT,                             -- Optional field for item description
    order_status TEXT DEFAULT 'Pending',               -- Default status of the order
    created_at TIMESTAMP DEFAULT NOW(),                -- Timestamp when the order was created
    updated_at TIMESTAMP DEFAULT NOW(),                -- Timestamp for the last update
    order_type_id INT NOT NULL,                        -- Reference to the order type
    total_fee FLOAT,                                   -- Total fee for delivery
    cod_fee FLOAT,                                     -- Cash on Delivery fee
    promo_discount FLOAT,                              -- Discount via promo codes
    discount FLOAT,                                    -- Manual or offer discount
    delivery_fee FLOAT,                                -- Delivery fee
    archive Boolean, 
    FOREIGN KEY (userId) REFERENCES users (id)
);
