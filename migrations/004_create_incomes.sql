CREATE TABLE IF NOT EXISTS incomes (
    id                   BIGINT        NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name                 VARCHAR(255)  NOT NULL,
    amount               DECIMAL(18,2) NOT NULL,
    date                 DATE          NOT NULL,
    income_category_id   BIGINT        NOT NULL,
    account_id           BIGINT        NOT NULL,
    created_at           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_incomes_category FOREIGN KEY (income_category_id) REFERENCES income_categories(id),
    CONSTRAINT fk_incomes_account  FOREIGN KEY (account_id)         REFERENCES accounts(id)
);
