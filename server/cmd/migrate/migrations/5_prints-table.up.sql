CREATE TABLE prints (
    `id` BINARY(16) PRIMARY KEY,
    `filament_id` BINARY(16) NOT NULL,
    `path_to_file` TEXT NOT NULL,
    `filament_amount_used` INTEGER NOT NULL,
    `estimated_print_time` TIME,

    FOREIGN KEY (filament_id) REFERENCES filaments(id)
);
