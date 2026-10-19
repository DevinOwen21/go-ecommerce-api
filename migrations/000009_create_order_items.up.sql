CREATE TABLE `order_items`
(
    id int NOT NULL AUTO_INCREMENT,
    order_id int NOT NULL,
    product_id int NOT NULL,
    quantity int NOT NULL,
    unit_price int NOT NULL,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`order_id`) REFERENCES orders(id) ON DELETE CASCADE ,
    FOREIGN KEY (`product_id`) REFERENCES products(id) ON DELETE RESTRICT,
    CONSTRAINT UC_order_items UNIQUE (order_id, product_id),
    CONSTRAINT chk_quantity CHECK (quantity > 0),
    CONSTRAINT chk_unit_price CHECK (unit_price >= 0)
)