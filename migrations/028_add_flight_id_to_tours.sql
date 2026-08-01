ALTER TABLE tours
  ADD COLUMN flight_id BIGINT DEFAULT NULL,
  ADD CONSTRAINT fk_tour_flight FOREIGN KEY (flight_id) REFERENCES flights(id);
