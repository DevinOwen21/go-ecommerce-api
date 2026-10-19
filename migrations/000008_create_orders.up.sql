CREATE TABLE `orders`
(
    id int NOT NULL AUTO_INCREMENT,
    user_id int NOT NULL,
    total_amount int NOT NULL,
    status varchar(20) NOT NULL DEFAULT 'pending',
    created_at datetime DEFAULT CURRENT_TIMESTAMP,
    expires_at datetime NOT NULL,

    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT chk_order_total CHECK (total_amount >= 0),
    CONSTRAINT chk_order_status CHECK (status IN ('pending', 'paid', 'cancelled', 'expired'))
)
