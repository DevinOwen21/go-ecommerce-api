ALTER TABLE products
    ADD COLUMN reserved_stock INT NOT NULL DEFAULT 0,
    ADD CONSTRAINT chk_reserved_stock
        CHECK (reserved_stock >= 0 AND reserved_stock <= stock);