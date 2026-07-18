ALTER TABLE clients
    ADD COLUMN birth_date DATE NULL AFTER phone;

UPDATE clients SET birth_date = DATE(CONCAT(birth_year, '-01-01')) WHERE birth_date IS NULL;

ALTER TABLE clients
    MODIFY COLUMN birth_date DATE NOT NULL,
    DROP COLUMN birth_year;
