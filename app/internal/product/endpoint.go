package product

import (
	"context"
	"fmt"

	"errors"
	"log"

	"github.com/ncostamagna/go-http-utils/meta"
	"github.com/ncostamagna/go-http-utils/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	// Endpoints struct
	Endpoints struct {
		Get    Controller
		GetAll Controller
		Store  Controller
		Update Controller
		Delete Controller
	}

	StoreReq struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	GetReq struct {
		ID string `json:"productId"`
	}

	GetAllReq struct {
		Name  string
		Limit int
		Page  int
	}

	UpdateReq struct {
		ID string
	}

	DeleteReq struct {
		ID string
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(s Service, c Config) Endpoints {
	return Endpoints{
		Get:    makeGet(s),
		GetAll: makeGetAll(s, c),
		Store:  makeStore(s),
		Update: makeUpdate(s),
		Delete: makeDelete(s),
	}
}

func makeGet(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)
		log.Println("Entra Get")
		log.Println(req)
		product, err := service.Get(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("Success", product, nil), nil
	}
}

func makeGetAll(service Service, c Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllReq)
		fmt.Println(req)
		filters := Filters{
			Name: req.Name,
		}

		count, err := service.Count(ctx, filters)
		if err != nil {
			fmt.Println(err)
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, c.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		products, err := service.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("Success", products, meta), nil
	}
}

func makeStore(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreReq)

		if req.Name == "" {
			return nil, response.BadRequest("Name is required")
		}

		if req.Price == 0 {
			return nil, response.BadRequest("Price is required")
		}

		product, err := service.Store(ctx, req.Name, req.Description, req.Price)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.Created("Success", product, nil), nil
	}
}

func makeUpdate(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return response.OK("Success", "UPDATE: testing 1234 6789", nil), nil
	}
}

func makeDelete(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.Println("Entra Delete")
		return response.OK("Success", "DELETE: testing 1234 6789", nil), nil
	}
}
