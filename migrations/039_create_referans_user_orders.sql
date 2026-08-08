-- Связь «реферер → заказ»: подтверждённые рефералы.
-- Кандидаты подбираются по clients.reference_name, а здесь хранится ручное подтверждение.

CREATE TABLE IF NOT EXISTS referans_user_orders (
    id               BIGINT   NOT NULL AUTO_INCREMENT,
    referans_user_id BIGINT   NOT NULL,
    order_id         BIGINT   NOT NULL,
    created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uniq_referans_user_order (referans_user_id, order_id),
    KEY idx_referans_user_orders_order (order_id),
    CONSTRAINT fk_referans_user_orders_user  FOREIGN KEY (referans_user_id) REFERENCES referans_users(id) ON DELETE CASCADE,
    CONSTRAINT fk_referans_user_orders_order FOREIGN KEY (order_id)         REFERENCES orders(id)         ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
