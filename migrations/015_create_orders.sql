CREATE TABLE IF NOT EXISTS orders (
    id         BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    client_id  BIGINT       NOT NULL,
    tour_id    BIGINT       NOT NULL,
    income_id  BIGINT       DEFAULT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_orders_client  FOREIGN KEY (client_id)  REFERENCES clients(id),
    CONSTRAINT fk_orders_tour    FOREIGN KEY (tour_id)    REFERENCES tours(id),
    CONSTRAINT fk_orders_income  FOREIGN KEY (income_id)  REFERENCES incomes(id)
);
