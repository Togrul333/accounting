ALTER TABLE clients
    ADD COLUMN gender        VARCHAR(10)  NOT NULL DEFAULT '' AFTER birth_date,
    ADD COLUMN nationality   VARCHAR(50)  NOT NULL DEFAULT '' AFTER gender,
    ADD COLUMN father_name   VARCHAR(100) NOT NULL DEFAULT '' AFTER nationality,
    ADD COLUMN reference_name VARCHAR(100) NOT NULL DEFAULT '' AFTER father_name;
