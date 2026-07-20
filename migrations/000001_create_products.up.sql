CREATE TABLE `products`
(
    `id`          int            NOT NULL AUTO_INCREMENT,
    `name`        varchar(100)   NOT NULL,
    `description` text,
    `price`       decimal(10, 2) NOT NULL,
    `stock`       int            NOT NULL,
    `created_at`  datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)