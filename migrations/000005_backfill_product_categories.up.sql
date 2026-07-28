UPDATE products
SET category = 1
WHERE category IS NULL;

ALTER TABLE products
    MODIFY COLUMN category INT NOT NULL;