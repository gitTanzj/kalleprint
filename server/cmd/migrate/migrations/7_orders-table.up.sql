CREATE TABLE orders (
    `id` BINARY(16) PRIMARY KEY,
    `customer_id` BINARY(16) NOT NULL,
    `print_id` BINARY(16) NOT NULL,
    `total` FLOAT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (`customer_id`) REFERENCES customers(`id`)
);
