CREATE TABLE `users`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(50) NOT NULL,
    `email` varchar(50) NOT NULL,
    `password_hash` varchar(255) NOT NULL,
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `role` enum('admin','customer') NOT NULL DEFAULT 'customer',
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`)
)