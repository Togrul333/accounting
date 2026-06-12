ALTER TABLE discounts ADD COLUMN order_id BIGINT DEFAULT NULL AFTER discount_category_id;
ALTER TABLE discounts ADD CONSTRAINT fk_discounts_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL;
