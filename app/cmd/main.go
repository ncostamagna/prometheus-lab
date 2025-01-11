package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/ncostamagna/prometheus-lab/app/internal/product"
	"github.com/ncostamagna/prometheus-lab/app/pkg/handler"

	"context"
	"flag"
	"log"
	"net/http"
	"os"
)

const writeTimeout = 10 * time.Second
const readTimeout = 4 * time.Second
const defaultURL = "0.0.0.0:80"

func main() {

	defer func() {
        if r := recover(); r != nil {
            log.Printf("Application panicked: %v", r)
        }
    }()

	//logger := bootstrap.NewLogger()

	_ = godotenv.Load()

	//logger.Info("DataBases")
	//db, err := bootstrap.DBConnection()
	//if err != nil {
	//	logger.Error(err)
	//	os.Exit(-1)
	//}

	flag.Parse()
	ctx := context.Background()

	
	var service product.Service
	{
		repository := product.NewRepo(nil, nil)
		service = product.NewService(nil, repository)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		os.Exit(-1)
	}

	h := handler.NewHTTPServer(ctx, product.MakeEndpoints(service, product.Config{LimPageDef: pagLimDef}))

	url := os.Getenv("APP_URL")
	if url == "" {
		url = defaultURL
	}

	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         url,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	errs := make(chan error)

	go func() {
		log.Println("listening on " + url)
		errs <- srv.ListenAndServe()
	}()

	err := <-errs
	if err != nil {
		log.Println(err)
	}

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == http.MethodOptions {
			return
		}

		h.ServeHTTP(w, r)
	})
}
