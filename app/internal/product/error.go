package product

import (
	"errors"
	"fmt"
)

var ErrNameRequired = errors.New("name is required")

type ErrNotFound struct {
	ProductID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("product '%s' doesn't exist", e.ProductID)
}
