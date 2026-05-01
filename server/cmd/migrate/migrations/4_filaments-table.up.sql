CREATE TABLE filaments (
    `id` BINARY(16) PRIMARY KEY,
    `filament_name` VARCHAR(100) NOT NULL,
    `current_amount` FLOAT NOT NULL,
    `cost_per_gram` FLOAT NOT NULL DEFAULT 0
);
