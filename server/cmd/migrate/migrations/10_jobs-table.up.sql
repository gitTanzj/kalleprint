CREATE TABLE jobs (
    `id` BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    `order_id` BINARY(16) NOT NULL,
    `filament_id` BINARY(16) NOT NULL,
    `gcode_path` TEXT NOT NULL,
    `status` ENUM('pending', 'printing', 'done') NOT NULL DEFAULT 'pending',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`order_id`) REFERENCES orders(`id`),
    FOREIGN KEY (`filament_id`) REFERENCES filaments(`id`)
);
