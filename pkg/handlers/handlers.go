package handlers

import (
	"github.com/RakhimovAns/Shop/cmd/server/app"
	"github.com/RakhimovAns/Shop/pkg/service"
)

type Server struct {
	server *app.Server
}

func NewServer(server *app.Server) *Server {
	return &Server{server: server}
}

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

var channel = make(chan *int64, 4)

func (s *Server) Init() {
	SubRoutineCustomers := s.server.Mux.PathPrefix("/api/customer").Subrouter()
	SubRoutineCustomers.HandleFunc("/save", s.HandleRegister).Methods(POST)
	SubRoutineCustomers.HandleFunc("/login", s.HandleLogin).Methods(GET)

	SubRoutineCustomer := s.server.Mux.PathPrefix("/api/customer").Subrouter()
	SubRoutineCustomer.Use(service.Auth(channel))
	SubRoutineCustomer.HandleFunc("/MakePurchase", s.HandleMakePurchase)
	SubRoutineCustomer.HandleFunc("/delete", s.HandleDelete).Methods(DELETE)
	SubRoutineCustomer.HandleFunc("/deposit", s.HandleDepositBalance).Methods(POST)
	SubRoutineCustomer.HandleFunc("/get/{ID}", s.HandleGetByID).Methods(GET)

	SubRoutineProduct := s.server.Mux.PathPrefix("/api/products").Subrouter()
	SubRoutineProduct.HandleFunc("", s.HandleGetProduct).Methods(GET)
	SubRoutineProduct.HandleFunc("/category", s.HandleGetCategory).Methods(GET)
	SubRoutineProduct.HandleFunc("/GetBy/{category}", s.HandleGetProductByCategory).Methods(GET)
	SubRoutineProduct.HandleFunc("/search/{Word}", s.HandleGetProductBySearch).Methods(GET)

	SubRoutineCart := s.server.Mux.PathPrefix("/api/cart").Subrouter()
	SubRoutineCart.Use(service.Auth(channel))
	SubRoutineCart.HandleFunc("/buy/products", s.HandleAddToCart).Methods(POST)
	SubRoutineCart.HandleFunc("/delete/products", s.HandleDeleteProduct).Methods(DELETE)
	SubRoutineCart.HandleFunc("/change/products", s.HandleChangeQTY).Methods(POST)
	SubRoutineCart.HandleFunc("/Get", s.HandleGetCart).Methods(GET)
	SubRoutineCart.HandleFunc("/Sum", s.HandleGetSum).Methods(GET)
	SubRoutineCart.HandleFunc("/Delete", s.HandleDeleteCart).Methods(DELETE)

	SubRoutinePurchase := s.server.Mux.PathPrefix("/api/purchase").Subrouter()
	SubRoutinePurchase.Use(service.Auth(channel))
	SubRoutinePurchase.HandleFunc("/GetAll", s.HandleGetAllPurchaseByID).Methods(GET)
}
