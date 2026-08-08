package repository

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type ReferansUserRepository interface {
	GetAll(ctx context.Context) ([]model.ReferansUser, error)
	GetByID(ctx context.Context, id int64) (*model.ReferansUser, error)
	Create(ctx context.Context, req model.CreateReferansUserRequest) (*model.ReferansUser, error)
	Update(ctx context.Context, id int64, req model.UpdateReferansUserRequest) (*model.ReferansUser, error)
	Delete(ctx context.Context, id int64) error

	// Candidates — заказы, где clients.reference_name похож на имя или фамилию реферера.
	Candidates(ctx context.Context, u model.ReferansUser, limit int) ([]model.ReferansOrder, error)
	// SearchOrders — поиск по всем заказам (клиент, референс, код тура, номер заказа).
	SearchOrders(ctx context.Context, userID int64, query string, limit int) ([]model.ReferansOrder, error)
	// Referrals — подтверждённые рефералы.
	Referrals(ctx context.Context, userID int64) ([]model.ReferansOrder, error)
	AddReferral(ctx context.Context, userID, orderID int64) error
	RemoveReferral(ctx context.Context, userID, orderID int64) error
}

type referansUserRepo struct {
	db *gorm.DB
}

func NewReferansUserRepository(db *gorm.DB) ReferansUserRepository {
	return &referansUserRepo{db: db}
}

func (r *referansUserRepo) GetAll(ctx context.Context) ([]model.ReferansUser, error) {
	var users []model.ReferansUser
	err := r.db.WithContext(ctx).Order("id").Find(&users).Error
	return users, err
}

func (r *referansUserRepo) GetByID(ctx context.Context, id int64) (*model.ReferansUser, error) {
	var u model.ReferansUser
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *referansUserRepo) Create(ctx context.Context, req model.CreateReferansUserRequest) (*model.ReferansUser, error) {
	u := model.ReferansUser{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}
	if err := r.db.WithContext(ctx).Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *referansUserRepo) Update(ctx context.Context, id int64, req model.UpdateReferansUserRequest) (*model.ReferansUser, error) {
	result := r.db.WithContext(ctx).Model(&model.ReferansUser{}).Where("id = ?", id).Updates(map[string]any{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"phone":      req.Phone,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		// Данные могли не измениться — проверяем, существует ли запись.
		if _, err := r.GetByID(ctx, id); err != nil {
			return nil, err
		}
	}
	return r.GetByID(ctx, id)
}

func (r *referansUserRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.ReferansUser{}, id).Error
}

// referansOrderQuery — общий SELECT для списков заказов реферера.
// Первый плейсхолдер (?) — id реферера, по нему считается признак confirmed.
const referansOrderQuery = `
	SELECT o.id, o.client_id,
	       CONCAT(c.first_name, ' ', c.last_name) AS client_name,
	       c.reference_name,
	       t.code AS tour_code,
	       COALESCE(inc.income_total, 0) AS income_total,
	       (ruo.id IS NOT NULL) AS confirmed,
	       o.created_at
	FROM orders o
	JOIN clients c ON c.id = o.client_id
	JOIN tours t   ON t.id = o.tour_id
	LEFT JOIN (
	    SELECT order_id, SUM(amount) AS income_total
	    FROM incomes GROUP BY order_id
	) inc ON inc.order_id = o.id
	LEFT JOIN referans_user_orders ruo ON ruo.order_id = o.id AND ruo.referans_user_id = ?`

func (r *referansUserRepo) Candidates(ctx context.Context, u model.ReferansUser, limit int) ([]model.ReferansOrder, error) {
	// Совпадение по имени или фамилии — сравниваем через LIKE,
	// так как reference_name у клиента обычно записан целиком («ADI SOYADI»).
	var conds []string
	args := []any{u.ID}
	for _, part := range []string{u.FirstName, u.LastName} {
		part = strings.TrimSpace(part)
		if len([]rune(part)) < 3 {
			continue
		}
		conds = append(conds, "c.reference_name LIKE ?")
		args = append(args, "%"+part+"%")
	}
	if len(conds) == 0 {
		return []model.ReferansOrder{}, nil
	}

	query := referansOrderQuery + " WHERE c.reference_name <> '' AND (" + strings.Join(conds, " OR ") + ") ORDER BY o.id DESC"
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	var orders []model.ReferansOrder
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&orders).Error
	if orders == nil {
		orders = []model.ReferansOrder{}
	}
	return orders, err
}

func (r *referansUserRepo) SearchOrders(ctx context.Context, userID int64, query string, limit int) ([]model.ReferansOrder, error) {
	query = strings.TrimSpace(query)
	args := []any{userID}
	sql := referansOrderQuery
	if query != "" {
		like := "%" + query + "%"
		sql += ` WHERE CONCAT(c.first_name, ' ', c.last_name) LIKE ?
		           OR c.reference_name LIKE ?
		           OR c.phone LIKE ?
		           OR t.code LIKE ?
		           OR CAST(o.id AS CHAR) = ?`
		args = append(args, like, like, like, like, query)
	}
	sql += " ORDER BY o.id DESC"
	if limit > 0 {
		sql += " LIMIT ?"
		args = append(args, limit)
	}

	var orders []model.ReferansOrder
	err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&orders).Error
	if orders == nil {
		orders = []model.ReferansOrder{}
	}
	return orders, err
}

func (r *referansUserRepo) Referrals(ctx context.Context, userID int64) ([]model.ReferansOrder, error) {
	sql := referansOrderQuery + " WHERE ruo.id IS NOT NULL ORDER BY o.id DESC"
	var orders []model.ReferansOrder
	err := r.db.WithContext(ctx).Raw(sql, userID).Scan(&orders).Error
	if orders == nil {
		orders = []model.ReferansOrder{}
	}
	return orders, err
}

func (r *referansUserRepo) AddReferral(ctx context.Context, userID, orderID int64) error {
	// INSERT IGNORE проглотил бы и ошибку внешнего ключа, поэтому наличие заказа
	// проверяем отдельно, а повторное подтверждение делаем идемпотентным.
	var cnt int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM orders WHERE id = ?", orderID).Scan(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO referans_user_orders (referans_user_id, order_id) VALUES (?, ?) "+
			"ON DUPLICATE KEY UPDATE referans_user_id = referans_user_id",
		userID, orderID,
	).Error
}

func (r *referansUserRepo) RemoveReferral(ctx context.Context, userID, orderID int64) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM referans_user_orders WHERE referans_user_id = ? AND order_id = ?",
		userID, orderID,
	).Error
}
