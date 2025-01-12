package handler

import (
	"github.com/gin-gonic/gin"
		"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/ncostamagna/go-http-utils/response"
	"github.com/ncostamagna/prometheus-lab/app/internal/product"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"context"
	"net/http"
	"encoding/json"
)

type ctxKey string

const (
	ctxParam  ctxKey = "params"
	ctxHeader ctxKey = "header"
)

func prometheusHandler() gin.HandlerFunc {
    h := promhttp.Handler()

    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}

func NewHTTPServer(_ context.Context, endpoints product.Endpoints) http.Handler {

	r := gin.Default()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}
	r.Use(ginDecode())

	r.GET("/products", gin.WrapH(httptransport.NewServer(endpoint.Endpoint(endpoints.Get), decodeGetHandler, encodeResponse, opts...)))
	r.POST("/products", gin.WrapH(httptransport.NewServer(endpoint.Endpoint(endpoints.Store), decodeStoreHandler, encodeResponse, opts...)))

	r.GET("/metrics", prometheusHandler())

	return r
}

func ginDecode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ctxParam, c.Params)
		ctx = context.WithValue(ctx, ctxHeader, c.Request.Header)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func decodeGetHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	pp := ctx.Value(ctxHeader).(http.Header)

	if len(pp["Authorization"]) < 1 {
		return nil, response.BadRequest("invalid authentication")
	}

	req := product.GetReq{}

	return req, nil
}

func decodeGetAllHandler(_ context.Context,  r *http.Request) (interface{}, error) {

	return nil, nil
}

func decodeStoreHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	var req product.StoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateHandler(ctx context.Context, r *http.Request) (interface{}, error) {

	var req product.UpdateReq

	return req, nil
}

func decodeDeleteHandler(ctx context.Context, r *http.Request) (interface{}, error) {

	var req product.DeleteReq

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
