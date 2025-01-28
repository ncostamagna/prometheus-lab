package product

import (
	"cmp"
	"context"
	"slices"

	"github.com/ncostamagna/go-logger-hub/loghub"

	"github.com/ncostamagna/prometheus-lab/app/internal/domain"
)

type (
	Repository interface {
		Store(ctx context.Context, product *domain.Product) error
		GetAll(ctx context.Context, offset, limit int) ([]domain.Product, error)
		Get(ctx context.Context, id int) (*domain.Product, error)
		Delete(ctx context.Context, id int) error
		Update(ctx context.Context, id int, name, description *string, price *float64) error
		Count(ctx context.Context) (int, error)
	}

	db struct {
		products []domain.Product
		maxID    int
	}
	repo struct {
		db  db
		log loghub.Logger
	}
)

// NewRepo is a repositories handler.
func NewRepo(l loghub.Logger) Repository {
	return &repo{
		db: db{
			products: []domain.Product{},
			maxID:    0,
		},
		log: l,
	}
}

func (r *repo) Store(_ context.Context, product *domain.Product) error {
	r.db.maxID++
	product.ID = r.db.maxID
	r.db.products = append(r.db.products, *product)
	return nil
}

func (r *repo) GetAll(_ context.Context, _, _ int) ([]domain.Product, error) {
	return r.db.products, nil
}

func (r *repo) Get(_ context.Context, id int) (*domain.Product, error) {

	i, found := slices.BinarySearchFunc(r.db.products, id, func(a domain.Product, b int) int {
		return cmp.Compare(a.ID, b)
	})

	if !found {
		return nil, ErrNotFound{id}
	}

	return &r.db.products[i], nil
}

func (r *repo) Delete(_ context.Context, id int) error {
	r.db.products = slices.DeleteFunc(r.db.products, func(a domain.Product) bool {
		return a.ID == id
	})
	return nil
}

func (r *repo) Update(ctx context.Context, id int, name, description *string, price *float64) error {
	p, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	if name != nil {
		p.Name = *name
	}

	if description != nil {
		p.Description = *description
	}

	if price != nil {
		p.Price = *price
	}

	return nil
}

func (r *repo) Count(_ context.Context) (int, error) {
	return r.db.maxID, nil
}
