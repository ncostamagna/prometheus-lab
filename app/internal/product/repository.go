package product

import (
	"cmp"
	"context"
	"github.com/ncostamagna/prometheus-lab/app/internal/domain"
	"gorm.io/gorm"
	"log"
	"slices"
)

var products []domain.Product
var maxId int

func init() {
	products = []domain.Product{}
}

type (
	Repository interface {
		Store(ctx context.Context, product *domain.Product) error
		GetAll(ctx context.Context, offset, limit int) ([]domain.Product, error)
		Get(ctx context.Context, id int) (*domain.Product, error)
		Delete(ctx context.Context, id int) error
		Update(ctx context.Context, id string, name, description *string, price *float64) error
		Count(ctx context.Context) (int, error)
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
	maxId++
	product.ID = maxId
	products = append(products, *product)
	return nil
}

func (r *repo) GetAll(ctx context.Context, offset, limit int) ([]domain.Product, error) {
	return products, nil
}

func (r *repo) Get(ctx context.Context, id int) (*domain.Product, error) {

	i, found := slices.BinarySearchFunc(products, id, func(a domain.Product, b int) int {
		return cmp.Compare(a.ID, b)
	})

	if !found {
		return nil, ErrNotFound{id}
	}

	return &products[i], nil
}

func (r *repo) Delete(ctx context.Context, id int) error {
	slices.DeleteFunc(products, func(a domain.Product) bool {
		return a.ID == id
	})
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

func (r *repo) Count(ctx context.Context) (int, error) {
	return maxId, nil
}
