package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"accounting/internal/model"
)

type MetaAdAccountRepository interface {
	GetAll(ctx context.Context) ([]model.MetaAdAccount, error)
	GetByID(ctx context.Context, id int64) (*model.MetaAdAccount, error)
	Create(ctx context.Context, req model.CreateMetaAdAccountRequest) (*model.MetaAdAccount, error)
	Update(ctx context.Context, id int64, req model.UpdateMetaAdAccountRequest) (*model.MetaAdAccount, error)
	Delete(ctx context.Context, id int64) error
	TouchSynced(ctx context.Context, id int64, name, currency string) error
}

type MetaAdSpendRepository interface {
	Upsert(ctx context.Context, rows []model.MetaAdSpend) (int, error)
	List(ctx context.Context, adAccountID, since, until string) ([]model.MetaAdSpend, error)
	CampaignSummary(ctx context.Context, adAccountID, since, until string) ([]model.MetaCampaignSummary, error)
	Unbilled(ctx context.Context, adAccountID, since, until string) ([]model.MetaAdSpend, error)
	SetExpenseID(ctx context.Context, id, expenseID int64) error
	SetTourID(ctx context.Context, id int64, tourID *int64) error
}

// ── Рекламные кабинеты ───────────────────────────────────────────────────────

type metaAdAccountRepo struct {
	db *gorm.DB
}

func NewMetaAdAccountRepository(db *gorm.DB) MetaAdAccountRepository {
	return &metaAdAccountRepo{db: db}
}

const metaAdAccountSelectQuery = `
	SELECT m.id, m.ad_account_id, m.name, m.currency, m.access_token,
	       m.expense_category_id, c.name AS expense_category_name,
	       m.account_id, a.name AS account_name,
	       m.auto_expense, m.last_synced_at, m.created_at, m.updated_at
	FROM meta_ad_accounts m
	LEFT JOIN expense_categories c ON c.id = m.expense_category_id
	LEFT JOIN accounts a           ON a.id = m.account_id`

func (r *metaAdAccountRepo) GetAll(ctx context.Context) ([]model.MetaAdAccount, error) {
	var accounts []model.MetaAdAccount
	err := r.db.WithContext(ctx).Raw(metaAdAccountSelectQuery + ` ORDER BY m.id`).Scan(&accounts).Error
	if accounts == nil {
		accounts = []model.MetaAdAccount{}
	}
	return accounts, err
}

func (r *metaAdAccountRepo) GetByID(ctx context.Context, id int64) (*model.MetaAdAccount, error) {
	var acc model.MetaAdAccount
	result := r.db.WithContext(ctx).Raw(metaAdAccountSelectQuery+` WHERE m.id = ?`, id).Scan(&acc)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &acc, nil
}

func (r *metaAdAccountRepo) Create(ctx context.Context, req model.CreateMetaAdAccountRequest) (*model.MetaAdAccount, error) {
	acc := model.MetaAdAccount{
		AdAccountID:       req.AdAccountID,
		Name:              req.Name,
		AccessToken:       req.AccessToken,
		ExpenseCategoryID: req.ExpenseCategoryID,
		AccountID:         req.AccountID,
		AutoExpense:       req.AutoExpense,
	}
	if err := r.db.WithContext(ctx).Create(&acc).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, acc.ID)
}

func (r *metaAdAccountRepo) Update(ctx context.Context, id int64, req model.UpdateMetaAdAccountRequest) (*model.MetaAdAccount, error) {
	updates := map[string]any{
		"name":                req.Name,
		"expense_category_id": req.ExpenseCategoryID,
		"account_id":          req.AccountID,
		"auto_expense":        req.AutoExpense,
	}
	// Пустой токен в запросе означает «оставить прежний».
	if req.AccessToken != "" {
		updates["access_token"] = req.AccessToken
	}

	result := r.db.WithContext(ctx).Model(&model.MetaAdAccount{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		// Значение могло не измениться — проверяем, существует ли запись.
		if _, err := r.GetByID(ctx, id); err != nil {
			return nil, err
		}
	}
	return r.GetByID(ctx, id)
}

func (r *metaAdAccountRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.MetaAdAccount{}, id).Error
}

// TouchSynced обновляет время синхронизации и подтягивает имя/валюту из Meta.
func (r *metaAdAccountRepo) TouchSynced(ctx context.Context, id int64, name, currency string) error {
	updates := map[string]any{"last_synced_at": time.Now()}
	if name != "" {
		updates["name"] = name
	}
	if currency != "" {
		updates["currency"] = currency
	}
	return r.db.WithContext(ctx).Model(&model.MetaAdAccount{}).Where("id = ?", id).Updates(updates).Error
}

// ── Расходы по кампаниям ─────────────────────────────────────────────────────

type metaAdSpendRepo struct {
	db *gorm.DB
}

func NewMetaAdSpendRepository(db *gorm.DB) MetaAdSpendRepository {
	return &metaAdSpendRepo{db: db}
}

// Upsert пишет статистику по ключу (ad_account_id, campaign_id, date).
// expense_id и tour_id намеренно не перезаписываются — привязки живут дольше синка.
func (r *metaAdSpendRepo) Upsert(ctx context.Context, rows []model.MetaAdSpend) (int, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	const batchSize = 200
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "ad_account_id"}, {Name: "campaign_id"}, {Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"campaign_name", "spend", "impressions", "clicks", "reach", "cpc", "ctr", "currency",
			}),
		}).
		CreateInBatches(&rows, batchSize).Error
	if err != nil {
		return 0, err
	}
	return len(rows), nil
}

const metaSpendSelectQuery = `
	SELECT s.id, s.ad_account_id, s.campaign_id, s.campaign_name, s.date,
	       s.spend, s.impressions, s.clicks, s.reach, s.cpc, s.ctr, s.currency,
	       s.tour_id, s.expense_id, s.created_at, s.updated_at
	FROM meta_ad_spend s`

func (r *metaAdSpendRepo) List(ctx context.Context, adAccountID, since, until string) ([]model.MetaAdSpend, error) {
	var rows []model.MetaAdSpend
	err := r.db.WithContext(ctx).
		Raw(metaSpendSelectQuery+`
			WHERE s.ad_account_id = ? AND s.date BETWEEN ? AND ?
			ORDER BY s.date DESC, s.spend DESC`, adAccountID, since, until).
		Scan(&rows).Error
	if rows == nil {
		rows = []model.MetaAdSpend{}
	}
	return rows, err
}

func (r *metaAdSpendRepo) CampaignSummary(ctx context.Context, adAccountID, since, until string) ([]model.MetaCampaignSummary, error) {
	var rows []model.MetaCampaignSummary
	err := r.db.WithContext(ctx).Raw(`
		SELECT s.campaign_id,
		       MAX(s.campaign_name)      AS campaign_name,
		       SUM(s.spend)              AS spend,
		       SUM(s.impressions)        AS impressions,
		       SUM(s.clicks)             AS clicks,
		       MAX(s.currency)           AS currency,
		       COUNT(DISTINCT s.date)    AS days
		FROM meta_ad_spend s
		WHERE s.ad_account_id = ? AND s.date BETWEEN ? AND ?
		GROUP BY s.campaign_id
		ORDER BY spend DESC`, adAccountID, since, until).Scan(&rows).Error
	if rows == nil {
		rows = []model.MetaCampaignSummary{}
	}
	if err != nil {
		return rows, err
	}
	// CPC и CTR считаем от суммы, а не усреднением дневных значений.
	for i := range rows {
		if rows[i].Clicks > 0 {
			rows[i].CPC = rows[i].Spend / float64(rows[i].Clicks)
		}
		if rows[i].Impressions > 0 {
			rows[i].CTR = float64(rows[i].Clicks) / float64(rows[i].Impressions) * 100
		}
	}
	return rows, nil
}

// Unbilled — строки с расходом, для которых ещё не создана запись в expenses.
func (r *metaAdSpendRepo) Unbilled(ctx context.Context, adAccountID, since, until string) ([]model.MetaAdSpend, error) {
	var rows []model.MetaAdSpend
	err := r.db.WithContext(ctx).
		Raw(metaSpendSelectQuery+`
			WHERE s.ad_account_id = ? AND s.date BETWEEN ? AND ?
			  AND s.expense_id IS NULL AND s.spend > 0
			ORDER BY s.date`, adAccountID, since, until).
		Scan(&rows).Error
	if rows == nil {
		rows = []model.MetaAdSpend{}
	}
	return rows, err
}

func (r *metaAdSpendRepo) SetExpenseID(ctx context.Context, id, expenseID int64) error {
	return r.db.WithContext(ctx).Model(&model.MetaAdSpend{}).
		Where("id = ?", id).
		Update("expense_id", expenseID).Error
}

func (r *metaAdSpendRepo) SetTourID(ctx context.Context, id int64, tourID *int64) error {
	return r.db.WithContext(ctx).Model(&model.MetaAdSpend{}).
		Where("id = ?", id).
		Update("tour_id", tourID).Error
}
