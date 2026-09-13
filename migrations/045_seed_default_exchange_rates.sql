-- Döviz Kurları üçün ilkin default dəyərlər. Yalnız hələ heç kimin dəyişmədiyi
-- ('0' seed dəyərində qalan) sətirlər yenilənir ki, artıq əl ilə daxil edilmiş
-- kurslar üzərinə yazılmasın.
UPDATE settings SET value = '1.7000' WHERE `key` = 'rate_usd' AND value = '0';
UPDATE settings SET value = '0.0350' WHERE `key` = 'rate_try' AND value = '0';
UPDATE settings SET value = '1.9729' WHERE `key` = 'rate_eur' AND value = '0';
UPDATE settings SET value = '2.2956' WHERE `key` = 'rate_gbp' AND value = '0';
