CREATE TABLE orders (
    `id` BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    `name` VARCHAR(100) NOT NULL,
    `email` VARCHAR(100) NOT NULL,
    `shipping_address` TEXT NOT NULL,
    `notes` TEXT,
    `print_id` BINARY(16),
    `total` FLOAT NOT NULL,
    `status` ENUM('pending','printing','waiting-for-shipment','shipped','done') NOT NULL DEFAULT 'pending',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`)
);
