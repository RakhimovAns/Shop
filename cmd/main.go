package main

import (
	"context"
	"github.com/RakhimovAns/Shop/cmd/server/app"
	"github.com/RakhimovAns/Shop/pkg/carts"
	"github.com/RakhimovAns/Shop/pkg/customer"
	"github.com/RakhimovAns/Shop/pkg/products"
	"github.com/RakhimovAns/Shop/pkg/purchase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/dig"
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
	deps := []interface{}{
		app.NewServer,
		http.NewServeMux,
		mux.NewRouter,
		func() (*pgxpool.Pool, error) {
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(ctx, dsn)
		},
		customer.NewService,
		product.NewService,
		carts.NewService,
		purchase.NewService,
		func(server *app.Server) *http.Server {
			return &http.Server{
				Addr:    net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}

	container := dig.New()
	for _, dep := range deps {
		err := container.Provide(dep)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	err = container.Invoke(func(server *app.Server) {
		server.Init()
	})
	if err != nil {
		log.Println(err)
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}
