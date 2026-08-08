-- Пивот тур ↔ комнаты. Цена комнаты теперь хранится здесь, а не в rooms:
-- один и тот же тип номера может стоить по-разному в разных турах.
CREATE TABLE tour_rooms (
    tour_id BIGINT        NOT NULL,
    room_id BIGINT        NOT NULL,
    price   DECIMAL(12,2) NOT NULL DEFAULT 0,
    PRIMARY KEY (tour_id, room_id),
    CONSTRAINT fk_tour_rooms_tour FOREIGN KEY (tour_id) REFERENCES tours(id) ON DELETE CASCADE,
    CONSTRAINT fk_tour_rooms_room FOREIGN KEY (room_id) REFERENCES rooms(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
