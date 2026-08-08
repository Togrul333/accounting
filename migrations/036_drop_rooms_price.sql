-- rooms остаётся справочником типов номеров (код + количество мест).
-- Цена задаётся на туре, в tour_rooms.price.
ALTER TABLE rooms
  DROP COLUMN price;
