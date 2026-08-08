-- Meta (Facebook/Instagram) Ads интеграция: настройки рекламных кабинетов
-- и ежедневная статистика расходов по кампаниям.

CREATE TABLE IF NOT EXISTS meta_ad_accounts (
    id                  BIGINT       NOT NULL AUTO_INCREMENT,
    ad_account_id       VARCHAR(64)  NOT NULL,              -- act_1234567890
    name                VARCHAR(255) NOT NULL DEFAULT '',   -- название кабинета из Meta
    currency            VARCHAR(10)  NOT NULL DEFAULT '',   -- TRY / USD / ...
    access_token        TEXT         NULL,                  -- system user token; пусто = берём META_ACCESS_TOKEN из env
    expense_category_id BIGINT       NULL,                  -- в какую категорию писать расход
    account_id          BIGINT       NULL,                  -- с какого счёта списывать
    auto_expense        TINYINT(1)   NOT NULL DEFAULT 0,    -- создавать записи в expenses автоматически
    last_synced_at      DATETIME     NULL,
    created_at          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uniq_meta_ad_account (ad_account_id),
    CONSTRAINT fk_meta_ad_accounts_category FOREIGN KEY (expense_category_id) REFERENCES expense_categories(id) ON DELETE SET NULL,
    CONSTRAINT fk_meta_ad_accounts_account  FOREIGN KEY (account_id)          REFERENCES accounts(id)           ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS meta_ad_spend (
    id            BIGINT        NOT NULL AUTO_INCREMENT,
    ad_account_id VARCHAR(64)   NOT NULL,
    campaign_id   VARCHAR(64)   NOT NULL,
    campaign_name VARCHAR(500)  NOT NULL DEFAULT '',
    date          DATE          NOT NULL,
    spend         DECIMAL(18,2) NOT NULL DEFAULT 0,
    impressions   BIGINT        NOT NULL DEFAULT 0,
    clicks        BIGINT        NOT NULL DEFAULT 0,
    reach         BIGINT        NOT NULL DEFAULT 0,
    cpc           DECIMAL(18,4) NOT NULL DEFAULT 0,
    ctr           DECIMAL(10,4) NOT NULL DEFAULT 0,
    currency      VARCHAR(10)   NOT NULL DEFAULT '',
    tour_id       BIGINT        NULL,                       -- ручная привязка кампании к туру
    expense_id    BIGINT        NULL,                       -- созданная запись расхода (защита от дублей)
    created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uniq_meta_spend (ad_account_id, campaign_id, date),
    KEY idx_meta_spend_date (date),
    CONSTRAINT fk_meta_spend_tour    FOREIGN KEY (tour_id)    REFERENCES tours(id)    ON DELETE SET NULL,
    CONSTRAINT fk_meta_spend_expense FOREIGN KEY (expense_id) REFERENCES expenses(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Категория расходов для рекламы; name уникален, поэтому повторный запуск безопасен.
INSERT IGNORE INTO expense_categories (name) VALUES ('Reklam');
