CREATE TABLE tours (
    id               BIGINT       NOT NULL AUTO_INCREMENT,
    code             VARCHAR(100) NOT NULL,
    start_date       DATE         NOT NULL,
    end_date         DATE         NOT NULL,
    tour_category_id BIGINT       NOT NULL,
    room_id          BIGINT       NOT NULL,
    created_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_tour_category FOREIGN KEY (tour_category_id) REFERENCES tour_categories(id),
    CONSTRAINT fk_tour_room     FOREIGN KEY (room_id)          REFERENCES rooms(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
