-- Раньше импорт создавал отдельный тур на каждую пару (комната × категория),
-- поэтому в базе лежат дубли с одинаковым кодом. Схлопываем их в один тур,
-- а комнаты уводим в пивот.

-- 1. Переносим текущую привязку комнаты и её цену в пивот.
INSERT IGNORE INTO tour_rooms (tour_id, room_id, price)
SELECT t.id, t.room_id, r.price
  FROM tours t
  JOIN rooms r ON r.id = t.room_id
 WHERE t.room_id IS NOT NULL;

-- 2. Карта «дубль → выживший тур». Выживший — минимальный id в группе
--    (код, категория, даты). Обычная таблица, а не TEMPORARY: MySQL не даёт
--    сослаться на временную таблицу дважды в одном запросе.
DROP TABLE IF EXISTS tour_merge_map;
CREATE TABLE tour_merge_map (
    old_id BIGINT NOT NULL PRIMARY KEY,
    new_id BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO tour_merge_map (old_id, new_id)
SELECT t.id, s.keep_id
  FROM tours t
  JOIN (
        SELECT code, tour_category_id, start_date, end_date, MIN(id) AS keep_id
          FROM tours
         GROUP BY code, tour_category_id, start_date, end_date
       ) s
    ON s.code             = t.code
   AND s.tour_category_id = t.tour_category_id
   AND s.start_date       = t.start_date
   AND s.end_date         = t.end_date
 WHERE t.id <> s.keep_id;

-- 3. Связи комнат и рейсов переезжают на выжившего.
--    INSERT IGNORE гасит конфликты по составному PK, если у дублей были общие рейсы.
INSERT IGNORE INTO tour_rooms (tour_id, room_id, price)
SELECT m.new_id, tr.room_id, tr.price
  FROM tour_rooms tr
  JOIN tour_merge_map m ON m.old_id = tr.tour_id;

DELETE tr FROM tour_rooms tr
  JOIN tour_merge_map m ON m.old_id = tr.tour_id;

INSERT IGNORE INTO tour_flights (tour_id, flight_id)
SELECT m.new_id, tf.flight_id
  FROM tour_flights tf
  JOIN tour_merge_map m ON m.old_id = tf.tour_id;

DELETE tf FROM tour_flights tf
  JOIN tour_merge_map m ON m.old_id = tf.tour_id;

-- 4. Заказы, доходы и расходы перевешиваем на выжившего.
--    orders.room_id уже проставлен миграцией 033, поэтому цена номера не теряется.
UPDATE orders   o JOIN tour_merge_map m ON m.old_id = o.tour_id SET o.tour_id = m.new_id;
UPDATE incomes  i JOIN tour_merge_map m ON m.old_id = i.tour_id SET i.tour_id = m.new_id;
UPDATE expenses e JOIN tour_merge_map m ON m.old_id = e.tour_id SET e.tour_id = m.new_id;

-- 5. Дубли больше ни на что не ссылаются — удаляем.
DELETE t FROM tours t
  JOIN tour_merge_map m ON m.old_id = t.id;

DROP TABLE tour_merge_map;
