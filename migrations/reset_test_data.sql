-- Сброс тестовых данных перед реальным запуском.
-- Удаляет все туры, заказы, доходы и расходы, и всё, что зависит от них по FK.
--
-- Порядок важен: fk_orders_tour и fk_incomes_tour объявлены без ON DELETE,
-- поэтому orders и incomes надо удалить раньше tours, иначе MySQL откажет
-- с ошибкой "Cannot delete or update a parent row: a foreign key constraint fails".
--
-- Каскадом (ON DELETE CASCADE/SET NULL из миграций) само подчистится:
--   orders  -> referans_user_orders (CASCADE),
--             incomes.order_id / discounts.order_id / tasks.order_id (SET NULL)
--   tours   -> tasks с tour_id (CASCADE) -> их task_comments (CASCADE),
--             tour_flights, tour_rooms (CASCADE),
--             expenses.tour_id / meta_ad_spend.tour_id (SET NULL)
--
-- НЕ трогает: clients, accounts, users, категории (tour/income/expense/discount),
-- rooms, flights, sheet_links, meta_ad_accounts, referans_users, hoca_users,
-- default_tasks, settings.

START TRANSACTION;

DELETE FROM orders;
DELETE FROM incomes;
DELETE FROM expenses;
DELETE FROM tours;

COMMIT;

-- Автоинкременты — чтобы новые записи снова начинались с 1.
ALTER TABLE orders                AUTO_INCREMENT = 1;
ALTER TABLE incomes                AUTO_INCREMENT = 1;
ALTER TABLE expenses               AUTO_INCREMENT = 1;
ALTER TABLE tours                  AUTO_INCREMENT = 1;
ALTER TABLE tasks                  AUTO_INCREMENT = 1;
ALTER TABLE task_comments          AUTO_INCREMENT = 1;
ALTER TABLE referans_user_orders   AUTO_INCREMENT = 1;
