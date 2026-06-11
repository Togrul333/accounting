CREATE TABLE IF NOT EXISTS expenses (
    id                    BIGINT        NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name                  VARCHAR(255)  NOT NULL,
    amount                DECIMAL(18,2) NOT NULL,
    date                  DATE          NOT NULL,
    expense_category_id   BIGINT        NOT NULL,
    account_id            BIGINT        NOT NULL,
    created_at            DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_expenses_category FOREIGN KEY (expense_category_id) REFERENCES expense_categories(id),
    CONSTRAINT fk_expenses_account  FOREIGN KEY (account_id)          REFERENCES accounts(id)
);
