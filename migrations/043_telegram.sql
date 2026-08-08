-- Telegram-бот: настройки бота лежат в settings, привязка сотрудника — в hoca_users.
-- Бот не может написать человеку первым, поэтому сотрудник сам открывает бота
-- по ссылке t.me/<bot>?start=<code> и жмёт Start; после этого мы узнаём его chat_id.

INSERT INTO settings (`key`, `value`) VALUES
  ('telegram_bot_token', ''),
  ('telegram_bot_username', ''),
  ('telegram_enabled', '0'),
  ('telegram_update_offset', '0')
ON DUPLICATE KEY UPDATE `key` = `key`;

ALTER TABLE hoca_users ADD COLUMN telegram_chat_id   VARCHAR(64) NULL AFTER phone;
ALTER TABLE hoca_users ADD COLUMN telegram_username  VARCHAR(64) NULL AFTER telegram_chat_id;
ALTER TABLE hoca_users ADD COLUMN telegram_link_code VARCHAR(16) NULL AFTER telegram_username;

ALTER TABLE hoca_users ADD UNIQUE KEY uq_hoca_users_link_code (telegram_link_code);
