package product

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ncostamagna/prometheus-lab/app/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Store(ctx context.Context, product *domain.Product) error
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Product, error)
		Get(ctx context.Context, id string) (*domain.Product, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, name, description *string, price *float64) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

// NewRepo is a repositories handler
func NewRepo(db *gorm.DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Store(ctx context.Context, product *domain.Product) error {

	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		r.log.Println(err)
		return err
	}
	r.log.Println("course created with id: ", product.ID)
	return nil
}

func (r *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Product, error) {
	var p []domain.Product

	tx := r.db.WithContext(ctx).Model(&p)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&p)

	if result.Error != nil {
		r.log.Println(result.Error)
		return nil, result.Error
	}
	return p, nil
}

func (r *repo) Get(ctx context.Context, id string) (*domain.Product, error) {
	product := domain.Product{ID: id}

	if err := r.db.WithContext(ctx).First(&product).Error; err != nil {
		r.log.Println(err)
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound{id}
		}

		return nil, err
	}
	return &product, nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	product := domain.Product{ID: id}
	result := r.db.WithContext(ctx).Delete(&product)

	if result.Error != nil {
		r.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Printf("course %s doesn't exists", id)
		return ErrNotFound{id}
	}
	return nil
}

func (r *repo) Update(ctx context.Context, id string, name, description *string, price *float64) error {

	/*values := make(map[string]interface{})

	if name != nil {
		values["name"] = *name
	}

	if startDate != nil {
		values["start_date"] = *startDate
	}

	if endDate != nil {
		values["end_date"] = *endDate
	}

	result := r.db.WithContext(ctx).Model(&domain.Course{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		r.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Printf("course %s doesn't exists", id)
		return ErrNotFound{id}
	}*/

	return nil
}

func (r *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(domain.Product{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		r.log.Println(err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	return tx
}
