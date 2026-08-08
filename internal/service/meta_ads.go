package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"accounting/internal/meta"
	"accounting/internal/model"
	"accounting/internal/repository"
)

const dateLayout = "2006-01-02"

type MetaAdsService struct {
	accounts repository.MetaAdAccountRepository
	spend    repository.MetaAdSpendRepository
	expenses repository.ExpenseRepository
}

func NewMetaAdsService(
	accounts repository.MetaAdAccountRepository,
	spend repository.MetaAdSpendRepository,
	expenses repository.ExpenseRepository,
) *MetaAdsService {
	return &MetaAdsService{accounts: accounts, spend: spend, expenses: expenses}
}

// ── Кабинеты ─────────────────────────────────────────────────────────────────

func (s *MetaAdsService) GetAccounts(ctx context.Context) ([]model.MetaAdAccount, error) {
	accounts, err := s.accounts.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for i := range accounts {
		accounts[i].HasToken = s.tokenFor(&accounts[i]) != ""
	}
	return accounts, nil
}

func (s *MetaAdsService) GetAccount(ctx context.Context, id int64) (*model.MetaAdAccount, error) {
	acc, err := s.accounts.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	acc.HasToken = s.tokenFor(acc) != ""
	return acc, nil
}

func (s *MetaAdsService) CreateAccount(ctx context.Context, req model.CreateMetaAdAccountRequest) (*model.MetaAdAccount, error) {
	req.AdAccountID = meta.NormalizeAdAccountID(req.AdAccountID)
	if req.AdAccountID == "" {
		return nil, fmt.Errorf("reklam hesabı id-si tələb olunur")
	}
	return s.accounts.Create(ctx, req)
}

func (s *MetaAdsService) UpdateAccount(ctx context.Context, id int64, req model.UpdateMetaAdAccountRequest) (*model.MetaAdAccount, error) {
	return s.accounts.Update(ctx, id, req)
}

func (s *MetaAdsService) DeleteAccount(ctx context.Context, id int64) error {
	return s.accounts.Delete(ctx, id)
}

// tokenFor: токен из настроек кабинета, иначе — общий из окружения.
func (s *MetaAdsService) tokenFor(acc *model.MetaAdAccount) string {
	if acc != nil && acc.AccessToken != "" {
		return acc.AccessToken
	}
	return meta.EnvToken()
}

// Verify проверяет токен и доступ к кабинету, попутно сохраняя имя и валюту.
func (s *MetaAdsService) Verify(ctx context.Context, id int64) (*meta.AccountInfo, error) {
	acc, err := s.accounts.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	info, err := meta.NewClient(s.tokenFor(acc)).AccountInfo(ctx, acc.AdAccountID)
	if err != nil {
		return nil, err
	}
	if err := s.accounts.TouchSynced(ctx, id, info.Name, info.Currency); err != nil {
		log.Printf("meta ads: hesab məlumatı yenilənmədi: %v", err)
	}
	return info, nil
}

// ── Синхронизация ────────────────────────────────────────────────────────────

// Sync тянет дневную статистику по кампаниям за период и, если у кабинета
// включён auto_expense, создаёт соответствующие записи в expenses.
//
// Повторный запуск за тот же период безопасен: строки статистики обновляются
// по уникальному ключу, а расход создаётся только там, где expense_id пуст.
func (s *MetaAdsService) Sync(ctx context.Context, id int64, since, until string) (*model.MetaSyncResult, error) {
	acc, err := s.accounts.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	since, until, err = normalizeRange(since, until)
	if err != nil {
		return nil, err
	}

	token := s.tokenFor(acc)
	if token == "" {
		return nil, fmt.Errorf("access token yoxdur: kabinet ayarlarında və ya META_ACCESS_TOKEN-də təyin edin")
	}

	insights, err := meta.NewClient(token).CampaignInsights(ctx, acc.AdAccountID, since, until)
	if err != nil {
		return nil, err
	}

	result := &model.MetaSyncResult{
		AdAccountID: acc.AdAccountID,
		Since:       since,
		Until:       until,
		RowsFetched: len(insights),
		Currency:    acc.Currency,
	}

	rows := make([]model.MetaAdSpend, 0, len(insights))
	for _, in := range insights {
		date, err := time.Parse(dateLayout, in.Date)
		if err != nil {
			log.Printf("meta ads: tarix oxunmadı %q: %v", in.Date, err)
			continue
		}
		rows = append(rows, model.MetaAdSpend{
			AdAccountID:  acc.AdAccountID,
			CampaignID:   in.CampaignID,
			CampaignName: in.CampaignName,
			Date:         date,
			Spend:        in.Spend,
			Impressions:  in.Impressions,
			Clicks:       in.Clicks,
			Reach:        in.Reach,
			CPC:          in.CPC,
			CTR:          in.CTR,
			Currency:     in.Currency,
		})
		result.TotalSpend += in.Spend
		if in.Currency != "" {
			result.Currency = in.Currency
		}
	}

	saved, err := s.spend.Upsert(ctx, rows)
	if err != nil {
		return nil, err
	}
	result.RowsSaved = saved

	if acc.AutoExpense {
		created, err := s.createExpenses(ctx, acc, since, until)
		if err != nil {
			return nil, err
		}
		result.ExpensesCreated = created
	}

	if err := s.accounts.TouchSynced(ctx, id, "", result.Currency); err != nil {
		log.Printf("meta ads: son sinxronizasiya vaxtı yazılmadı: %v", err)
	}

	return result, nil
}

// createExpenses превращает несписанные строки статистики в записи расходов.
func (s *MetaAdsService) createExpenses(ctx context.Context, acc *model.MetaAdAccount, since, until string) (int, error) {
	if acc.ExpenseCategoryID == nil || acc.AccountID == nil {
		return 0, fmt.Errorf("avtomatik gider üçün gider kateqoriyası və hesab seçilməlidir")
	}

	pending, err := s.spend.Unbilled(ctx, acc.AdAccountID, since, until)
	if err != nil {
		return 0, err
	}

	created := 0
	for _, row := range pending {
		expense, err := s.expenses.Create(ctx, model.CreateExpenseRequest{
			Name:              fmt.Sprintf("Meta Reklam — %s", row.CampaignName),
			Amount:            row.Spend,
			Date:              row.Date.Format(dateLayout),
			ExpenseCategoryID: *acc.ExpenseCategoryID,
			AccountID:         *acc.AccountID,
			TourID:            row.TourID,
			// bank_ref используем как ключ идемпотентности — так же, как при импорте выписок.
			BankRef:      fmt.Sprintf("meta:%s:%s:%s", row.AdAccountID, row.CampaignID, row.Date.Format(dateLayout)),
			Counterparty: "Meta Platforms",
		})
		if err != nil {
			log.Printf("meta ads: gider yaradılmadı (campaign %s, %s): %v", row.CampaignID, row.Date.Format(dateLayout), err)
			continue
		}
		if err := s.spend.SetExpenseID(ctx, row.ID, expense.ID); err != nil {
			log.Printf("meta ads: expense_id yazılmadı (spend %d): %v", row.ID, err)
			continue
		}
		created++
	}
	return created, nil
}

// ── Чтение статистики ────────────────────────────────────────────────────────

func (s *MetaAdsService) Spend(ctx context.Context, adAccountID, since, until string) ([]model.MetaAdSpend, error) {
	since, until, err := normalizeRange(since, until)
	if err != nil {
		return nil, err
	}
	return s.spend.List(ctx, meta.NormalizeAdAccountID(adAccountID), since, until)
}

func (s *MetaAdsService) CampaignSummary(ctx context.Context, adAccountID, since, until string) ([]model.MetaCampaignSummary, error) {
	since, until, err := normalizeRange(since, until)
	if err != nil {
		return nil, err
	}
	return s.spend.CampaignSummary(ctx, meta.NormalizeAdAccountID(adAccountID), since, until)
}

func (s *MetaAdsService) SetTour(ctx context.Context, id int64, tourID *int64) error {
	return s.spend.SetTourID(ctx, id, tourID)
}

// normalizeRange проверяет даты и подставляет последние 30 дней, если период пуст.
func normalizeRange(since, until string) (string, string, error) {
	if until == "" {
		until = time.Now().Format(dateLayout)
	}
	if since == "" {
		since = time.Now().AddDate(0, 0, -29).Format(dateLayout)
	}

	from, err := time.Parse(dateLayout, since)
	if err != nil {
		return "", "", fmt.Errorf("başlanğıc tarixi keçərsizdir (YYYY-MM-DD)")
	}
	to, err := time.Parse(dateLayout, until)
	if err != nil {
		return "", "", fmt.Errorf("bitmə tarixi keçərsizdir (YYYY-MM-DD)")
	}
	if from.After(to) {
		return "", "", fmt.Errorf("başlanğıc tarixi bitmə tarixindən sonra ola bilməz")
	}
	return since, until, nil
}
