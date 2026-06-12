ALTER TABLE incomes
    ADD COLUMN tour_id BIGINT NULL AFTER account_id,
    ADD CONSTRAINT fk_incomes_tour FOREIGN KEY (tour_id) REFERENCES tours(id);
