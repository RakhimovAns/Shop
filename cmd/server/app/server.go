package app

import (
	"github.com/RakhimovAns/Shop/pkg/service"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	Mux          *mux.Router
	CustomersSvc *service.CustomerService
	ProductsSvc  *service.ProductService
	CartsSvc     *service.CartService
	PurchasesSvc *service.PurchaseService
}

// NewServer TODO Create NewServer
func NewServer(mux *mux.Router, customerSvc *service.CustomerService, productSvc *service.ProductService, cartsSvc *service.CartService, purchasesSvc *service.PurchaseService) *Server {
	return &Server{Mux: mux, CustomersSvc: customerSvc, ProductsSvc: productSvc, CartsSvc: cartsSvc, PurchasesSvc: purchasesSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.Mux.ServeHTTP(writer, request)
}
