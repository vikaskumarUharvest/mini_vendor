

-- Seed Data (1 user, 1 order, 1 item)

-- Insert User
INSERT INTO users (name, email, password, phone, age, city)
VALUES (
    'Vikas Prajapati',
    'vikas@example.com',
    'Mobile@2002',
    '9876543210',
    25,
    'Delhi'
);

-- Insert Order (status matches your Go: "created")
WITH inserted_order AS (
    INSERT INTO orders (amount, status)
    VALUES (1500.50, 'created')
    RETURNING id
)

-- Insert Order Item
INSERT INTO order_items (order_id, name, qty)
SELECT id, 'Laptop Bag', 1 FROM inserted_order;