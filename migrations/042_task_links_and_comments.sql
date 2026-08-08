-- Задачи можно связывать не только с туром, но и с заказом и клиентом.
-- Плюс приоритет и комментарии (как в Jira).

ALTER TABLE tasks ADD COLUMN priority  VARCHAR(20) NOT NULL DEFAULT 'medium' AFTER status;
ALTER TABLE tasks ADD COLUMN order_id  BIGINT NULL AFTER tour_id;
ALTER TABLE tasks ADD COLUMN client_id BIGINT NULL AFTER order_id;

ALTER TABLE tasks ADD CONSTRAINT fk_tasks_order
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL;
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_client
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE SET NULL;

-- Комментарии к задаче — лента обсуждения в попапе задачи.
CREATE TABLE IF NOT EXISTS task_comments (
    id         BIGINT   NOT NULL AUTO_INCREMENT,
    task_id    BIGINT   NOT NULL,
    user_id    BIGINT   NULL,          -- автор; NULL если пользователь удалён
    body       TEXT     NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_task_comments_task (task_id),
    CONSTRAINT fk_task_comments_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_task_comments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
