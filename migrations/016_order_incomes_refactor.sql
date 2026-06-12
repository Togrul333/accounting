-- убираем income_id с orders, добавляем order_id в incomes
ALTER TABLE orders DROP FOREIGN KEY fk_orders_income;
ALTER TABLE orders DROP COLUMN income_id;

ALTER TABLE incomes ADD COLUMN order_id BIGINT DEFAULT NULL AFTER tour_id;
ALTER TABLE incomes ADD CONSTRAINT fk_incomes_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL;
