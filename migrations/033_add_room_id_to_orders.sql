-- Комната переезжает с тура на заказ: у тура их теперь несколько,
-- а конкретный клиент бронирует один тип номера.
-- NULL допустим — заказы из банковской выписки создаются без комнаты.
ALTER TABLE orders
  ADD COLUMN room_id BIGINT DEFAULT NULL AFTER tour_id,
  ADD CONSTRAINT fk_orders_room FOREIGN KEY (room_id) REFERENCES rooms(id);

-- Бэкфилл, пока tours.room_id ещё существует.
UPDATE orders o
  JOIN tours t ON t.id = o.tour_id
   SET o.room_id = t.room_id;
