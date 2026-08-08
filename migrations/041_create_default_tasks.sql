-- Varsayılan görevler — шаблоны задач, которые автоматически создаются
-- при создании тура и привязываются к нему (tasks.tour_id).

CREATE TABLE IF NOT EXISTS default_tasks (
    id                BIGINT       NOT NULL AUTO_INCREMENT,
    title             VARCHAR(255) NOT NULL,
    description       TEXT         NULL,
    days_before_start INT          NOT NULL DEFAULT 0,  -- срок задачи = start_date - N дней
    hoca_user_id      BIGINT       NULL,                -- исполнитель по умолчанию
    created_at        DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_default_tasks_hoca_user FOREIGN KEY (hoca_user_id) REFERENCES hoca_users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE tasks ADD COLUMN tour_id BIGINT NULL AFTER hoca_user_id;
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_tour
    FOREIGN KEY (tour_id) REFERENCES tours(id) ON DELETE CASCADE;
