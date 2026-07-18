ALTER TABLE clients
    ADD COLUMN fin_code       VARCHAR(20)  NOT NULL DEFAULT '' AFTER birth_year,
    ADD COLUMN id_card_number VARCHAR(20)  NOT NULL DEFAULT '' AFTER fin_code,
    ADD COLUMN document_photo VARCHAR(255) NOT NULL DEFAULT '' AFTER id_card_number;
