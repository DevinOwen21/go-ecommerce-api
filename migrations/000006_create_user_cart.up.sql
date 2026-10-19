CREATE TABLE `user_carts`
(
    `user_id`          int            NOT NULL ,
    `product_id`        int   NOT NULL,
    `quantity` int NOT NULL,

    CHECK ( `quantity` > 0 ),

    PRIMARY KEY (`user_id`, `product_id`),

    FOREIGN KEY (`user_id`) REFERENCES users(id) ON DELETE CASCADE ,
    FOREIGN KEY (`product_id`) REFERENCES products(id) ON DELETE CASCADE
)