CREATE TABLE IF NOT EXISTS discounts (
  id                   BIGINT        NOT NULL AUTO_INCREMENT PRIMARY KEY,
  amount               DECIMAL(15,2) NOT NULL DEFAULT 0,
  discount_category_id BIGINT        NOT NULL,
  created_at           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_discounts_category FOREIGN KEY (discount_category_id) REFERENCES discount_categories(id)
);
