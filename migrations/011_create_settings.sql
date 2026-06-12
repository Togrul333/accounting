CREATE TABLE IF NOT EXISTS settings (
  `key`      VARCHAR(50)  NOT NULL PRIMARY KEY,
  `value`    VARCHAR(255) NOT NULL,
  updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO settings (`key`, `value`) VALUES
  ('rate_usd', '0'),
  ('rate_eur', '0'),
  ('rate_gbp', '0')
ON DUPLICATE KEY UPDATE `key` = `key`;
