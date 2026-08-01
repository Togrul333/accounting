INSERT INTO tour_flights (tour_id, flight_id)
SELECT id, flight_id FROM tours WHERE flight_id IS NOT NULL;

ALTER TABLE tours
  DROP FOREIGN KEY fk_tour_flight,
  DROP COLUMN flight_id;
