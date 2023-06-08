package app

import (
	"encoding/json"
	"github.com/RakhimovAns/Shop/cmd/app/middleware"
	"github.com/RakhimovAns/Shop/pkg/carts"
	"github.com/RakhimovAns/Shop/pkg/customer"
	product "github.com/RakhimovAns/Shop/pkg/products"
	"github.com/RakhimovAns/Shop/pkg/purchase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

type Server struct {
	mux          *mux.Router
	customersSvc *customer.Service
	productsSvc  *product.Service
	cartsSvc     *carts.Service
	purchasesSvc *purchase.Service
}

func NewServer(mux *mux.Router, customerSvc *customer.Service, productSvc *product.Service, cartsSvc *carts.Service, purchasesSvc *purchase.Service) *Server {
	return &Server{mux: mux, customersSvc: customerSvc, productsSvc: productSvc, cartsSvc: cartsSvc, purchasesSvc: purchasesSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

var channel = make(chan *int64, 4)

func (s *Server) Init() {
	SubRoutineProduct := s.mux.PathPrefix("/api/products").Subrouter()
	SubRoutineProduct.HandleFunc("", s.HandleGetProduct).Methods(GET)
	SubRoutineProduct.HandleFunc("/category", s.HandleGetCategory).Methods(GET)
	SubRoutineProduct.HandleFunc("/{category}", s.HandleGetProductByCategory).Methods(GET)
	SubRoutineProduct.HandleFunc("/search/{Word}", s.HandleGetProductBySearch).Methods(GET)

	SubRoutineCustomers := s.mux.PathPrefix("/api/customer").Subrouter()
	SubRoutineCustomers.HandleFunc("/save", s.HandleRegister).Methods(POST)
	SubRoutineCustomers.HandleFunc("/login", s.HandleLogin).Methods(POST)

	SubRoutineCustomer := s.mux.PathPrefix("/api/customer").Subrouter()
	SubRoutineCustomer.Use(middleware.Auth(channel))
	SubRoutineCustomer.HandleFunc("/MakePurchase", s.HandleMakePurchase)
	SubRoutineCustomer.HandleFunc("/check", s.HandleChecker).Methods(POST)
	SubRoutineCustomer.HandleFunc("/delete", s.HandleDelete).Methods(DELETE)
	SubRoutineCustomer.HandleFunc("/deposit", s.HandleDepositBalance).Methods(POST)
	s.mux.HandleFunc("/api/customer/get/{ID}", s.HandleGetByID).Methods(GET)

	SubRoutineCart := s.mux.PathPrefix("/api/cart").Subrouter()
	SubRoutineCart.Use(middleware.Auth(channel))
	SubRoutineCart.HandleFunc("/buy/products", s.HandleAddToCart).Methods(POST)
	SubRoutineCart.HandleFunc("/delete/products", s.HandleDeleteProduct).Methods(DELETE)
	SubRoutineCart.HandleFunc("/change/products", s.HandleChangeQTY).Methods(POST)
	SubRoutineCart.HandleFunc("/Get", s.HandleGetCart).Methods(GET)
	SubRoutineCart.HandleFunc("/Sum", s.HandleGetSum).Methods(GET)
	SubRoutineCart.HandleFunc("/Delete", s.HandleDeleteCart).Methods(DELETE)

	SubRoutinePurchase := s.mux.PathPrefix("/api/purchase").Subrouter()
	SubRoutinePurchase.Use(middleware.Auth(channel))
	SubRoutinePurchase.HandleFunc("/GetAll", s.HandleGetAllPurchaseByID).Methods(GET)

}
func (s *Server) HandleRegister(writer http.ResponseWriter, request *http.Request) {
	var item *customer.Customer
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.customersSvc.Register(request.Context(), item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = writer.Write([]byte("Was saved successfully"))
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
func (s *Server) HandleGetProduct(writer http.ResponseWriter, request *http.Request) {
	items, err := s.productsSvc.AllActiveProducts(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = writer.Write([]byte("There is all products\n"))
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (s *Server) HandleGetCategory(writer http.ResponseWriter, request *http.Request) {
	items, err := s.productsSvc.AllCategories(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}
	_, err = writer.Write([]byte("There is all categories\n"))
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (s *Server) HandleGetProductByCategory(writer http.ResponseWriter, request *http.Request) {
	Category := mux.Vars(request)["category"]
	items, err := s.productsSvc.GetByCategory(request.Context(), Category)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if len(items) == 0 {
		_, err = writer.Write([]byte("There is  not any product in this category"))
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = writer.Write([]byte("There is all products in this category\n"))
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (s *Server) HandleGetProductBySearch(writer http.ResponseWriter, request *http.Request) {
	Key := mux.Vars(request)["Word"]
	Key = "%" + Key + "%"
	items, err := s.productsSvc.Search(request.Context(), Key)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	if len(items) == 0 {
		_, err = writer.Write([]byte("There is  not any product for this key"))
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
func (s *Server) HandleLogin(writer http.ResponseWriter, request *http.Request) {
	var item *customer.Customer
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	token, err := s.customersSvc.Login(request.Context(), item.Phone, *item.Password)
	if err == customer.ErrNoSuchUser {
		_, err = writer.Write([]byte("No account with this phone number"))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	} else if err == customer.ErrInvalidPassword {
		_, err = writer.Write([]byte("Passwords don't match"))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	} else if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	} else if err == nil {
		_, err = writer.Write([]byte("You have been login successfully\nHere is your Token\n"))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		data, err := json.Marshal(token)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(data)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
}

func (s *Server) HandleDelete(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	err := s.customersSvc.Delete(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = writer.Write([]byte("Was deleted"))
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (s *Server) HandleGetByID(writer http.ResponseWriter, request *http.Request) {
	IDParam := mux.Vars(request)["ID"]
	id, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err, Customer := s.customersSvc.GetByID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(Customer)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (s *Server) HandleAddToCart(writer http.ResponseWriter, request *http.Request) {
	var Products *[]carts.Product
	err := json.NewDecoder(request.Body).Decode(&Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id := *<-channel
	err, ID := s.cartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.cartsSvc.SaveToCart(request.Context(), ID, Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = writer.Write([]byte("Added to cart successfully"))
}

func (s *Server) HandleDeleteProduct(writer http.ResponseWriter, request *http.Request) {
	var Products *[]carts.Product
	err := json.NewDecoder(request.Body).Decode(&Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id := *<-channel
	err, ID := s.cartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.cartsSvc.DeleteProducts(request.Context(), ID, Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (s *Server) HandleChangeQTY(writer http.ResponseWriter, request *http.Request) {
	var Products *[]carts.Product
	err := json.NewDecoder(request.Body).Decode(&Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id := *<-channel
	err, ID := s.cartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.cartsSvc.ChangeQTY(request.Context(), ID, Products)
	if err == carts.ErrNosuch {
		http.Error(writer, "You want to delete more than you have ", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (s *Server) HandleGetCart(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	err, ID := s.cartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	Products, err := s.cartsSvc.GetCartBYID(request.Context(), ID)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
func (s *Server) HandleGetSum(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	err, ID := s.cartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	Products, err := s.cartsSvc.GetCartBYID(request.Context(), ID)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	sum, err := s.cartsSvc.GetSum(request.Context(), Products)
	data, err := json.Marshal(sum)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
func (s *Server) HandleDeleteCart(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	err, ID := s.cartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.cartsSvc.DeleteCart(request.Context(), ID)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
func (s *Server) HandleMakePurchase(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	err, ID := s.cartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	Products, err := s.cartsSvc.GetCartBYID(request.Context(), ID)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	sum, err := s.cartsSvc.GetSum(request.Context(), Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err, Customer := s.customersSvc.GetByID(request.Context(), id)
	if Customer.Balance < sum {
		http.Error(writer, "You don't have enough money to buy it", http.StatusBadRequest)
		return
	}
	err = s.purchasesSvc.AddToPurchase(request.Context(), Products, id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.customersSvc.ChangeBalance(request.Context(), id, sum)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.cartsSvc.DeleteCart(request.Context(), ID)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (s *Server) HandleGetAllPurchaseByID(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	Purchases, err := s.purchasesSvc.GetAllPurchase(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(Purchases)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
func (s *Server) HandleDepositBalance(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	type Balance struct {
		Balance int64 `json:"balance"`
	}
	var item *Balance
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.customersSvc.DepositBalance(request.Context(), id, item.Balance)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
