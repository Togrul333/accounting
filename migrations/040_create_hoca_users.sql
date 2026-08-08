-- Hoca kullanıcıları — персонал турагентства (имя, фамилия, телефон).
-- На них можно назначать задачи (tasks.hoca_user_id).

CREATE TABLE IF NOT EXISTS hoca_users (
    id         BIGINT       NOT NULL AUTO_INCREMENT,
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name  VARCHAR(255) NOT NULL DEFAULT '',
    phone      VARCHAR(50)  NOT NULL DEFAULT '',
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE tasks ADD COLUMN hoca_user_id BIGINT NULL AFTER due_date;
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_hoca_user
    FOREIGN KEY (hoca_user_id) REFERENCES hoca_users(id) ON DELETE SET NULL;
