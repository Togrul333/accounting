-- Döviz Kurları artık AZN (manat) baz alınarak hesaplanıyor (əvvəl TRY idi).
-- TRY da digər valyutalar kimi manata çevrilməli olduğu üçün ayrıca kurs əlavə edilir.
INSERT INTO settings (`key`, `value`) VALUES
  ('rate_try', '0')
ON DUPLICATE KEY UPDATE `key` = `key`;
