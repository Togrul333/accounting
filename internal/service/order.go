package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type OrderService struct {
	repo       repository.OrderRepository
	incomeRepo repository.IncomeRepository
}

func NewOrderService(repo repository.OrderRepository, incomeRepo repository.IncomeRepository) *OrderService {
	return &OrderService{repo: repo, incomeRepo: incomeRepo}
}

func (s *OrderService) GetAll(ctx context.Context) ([]model.Order, error) {
	return s.repo.GetAll(ctx)
}

func (s *OrderService) GetByID(ctx context.Context, id int64) (*model.Order, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	incomes, err := s.incomeRepo.GetByOrderID(ctx, id)
	if err != nil {
		incomes = []model.Income{}
	}
	order.Incomes = incomes
	return order, nil
}

func (s *OrderService) Create(ctx context.Context, req model.CreateOrderRequest) (*model.Order, error) {
	order, err := s.repo.Create(ctx, req.ClientID, req.TourID)
	if err != nil {
		return nil, err
	}
	for _, incReq := range req.Incomes {
		incReq.OrderID = &order.ID
		if _, err := s.incomeRepo.Create(ctx, incReq); err != nil {
			return nil, err
		}
	}
	return s.GetByID(ctx, order.ID)
}

func (s *OrderService) AddIncome(ctx context.Context, orderID int64, req model.CreateIncomeRequest) (*model.Income, error) {
	req.OrderID = &orderID
	return s.incomeRepo.Create(ctx, req)
}

func (s *OrderService) Update(ctx context.Context, id int64, req model.UpdateOrderRequest) (*model.Order, error) {
	order, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	incomes, err := s.incomeRepo.GetByOrderID(ctx, id)
	if err != nil {
		incomes = []model.Income{}
	}
	order.Incomes = incomes
	return order, nil
}

func (s *OrderService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
