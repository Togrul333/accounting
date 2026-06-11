ALTER TABLE incomes
  ADD COLUMN account_id BIGINT NOT NULL AFTER income_category_id,
  ADD CONSTRAINT fk_incomes_account FOREIGN KEY (account_id) REFERENCES accounts(id);
