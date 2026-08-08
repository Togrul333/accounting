-- Связь тур → комната живёт в tour_rooms, колонка больше не нужна.
ALTER TABLE tours
  DROP FOREIGN KEY fk_tour_room,
  DROP COLUMN room_id;
