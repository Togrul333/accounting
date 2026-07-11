CREATE TABLE sheet_links (
    id             BIGINT        NOT NULL AUTO_INCREMENT,
    url            VARCHAR(1000) NOT NULL,
    spreadsheet_id VARCHAR(100)  NOT NULL,
    created_at     DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uniq_spreadsheet_id (spreadsheet_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
