package main

import (
	"context"
	"github.com/RakhimovAns/Shop/cmd/server/app"
	"github.com/RakhimovAns/Shop/pkg/handlers"
	"github.com/RakhimovAns/Shop/pkg/postgresql"
	"github.com/RakhimovAns/Shop/pkg/service"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgresql://postgres:Ansar@localhost:5432/db"
	if err := execute(host, port, dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(host string, port string, dsn string) (err error) {
	connectCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
	pool, err := pgxpool.Connect(connectCtx, dsn)
	if err != nil {
		log.Println(err)
		return
	}
	defer pool.Close()
	router := mux.NewRouter()
	CustomerSqlSvc := postgresql.NewCustomerService(pool)
	ProductSqlSvc := postgresql.NewProductService(pool)
	CartSqlSvc := postgresql.NewCartService(pool)
	PurchaseSqlSvc := postgresql.NewPurchaseService(pool)
	CustomerSvc := service.NewCustomerService(CustomerSqlSvc)
	ProductSvc := service.NewProductService(ProductSqlSvc)
	CartSvc := service.NewCartService(CartSqlSvc)
	PurchaseSvc := service.NewPurchaseService(PurchaseSqlSvc)
	server := app.NewServer(router, CustomerSvc, ProductSvc, CartSvc, PurchaseSvc)
	Server := handlers.NewServer(server)
	Server.Init()
	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}
	return srv.ListenAndServe()
}
