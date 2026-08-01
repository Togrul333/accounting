CREATE TABLE tour_flights (
    tour_id   BIGINT NOT NULL,
    flight_id BIGINT NOT NULL,
    PRIMARY KEY (tour_id, flight_id),
    CONSTRAINT fk_tour_flights_tour   FOREIGN KEY (tour_id)   REFERENCES tours(id) ON DELETE CASCADE,
    CONSTRAINT fk_tour_flights_flight FOREIGN KEY (flight_id) REFERENCES flights(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
