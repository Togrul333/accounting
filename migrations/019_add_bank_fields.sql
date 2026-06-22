ALTER TABLE incomes
  ADD COLUMN bank_ref            VARCHAR(100) DEFAULT NULL,
  ADD COLUMN counterparty        VARCHAR(500) DEFAULT NULL,
  ADD COLUMN counterparty_tax_id VARCHAR(50)  DEFAULT NULL;

ALTER TABLE expenses
  ADD COLUMN bank_ref            VARCHAR(100) DEFAULT NULL,
  ADD COLUMN counterparty        VARCHAR(500) DEFAULT NULL,
  ADD COLUMN counterparty_tax_id VARCHAR(50)  DEFAULT NULL;
