package app

import (
	"encoding/json"
	"github.com/RakhimovAns/Shop/cmd/app/middleware"
	"github.com/RakhimovAns/Shop/pkg/carts"
	"github.com/RakhimovAns/Shop/pkg/customer"
	product "github.com/RakhimovAns/Shop/pkg/products"
	"github.com/gorilla/mux"
	"log"
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
}

func NewServer(mux *mux.Router, customerSvc *customer.Service, productSvc *product.Service, cartsSvc *carts.Service) *Server {
	return &Server{mux: mux, customersSvc: customerSvc, productsSvc: productSvc, cartsSvc: cartsSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

var channel = make(chan *int64, 4)

func (s *Server) Init() {
	SubRoutineCustomer := s.mux.PathPrefix("/api/customer").Subrouter()
	SubRoutineCustomer.Use(middleware.Auth(channel))
	s.mux.HandleFunc("/api/customer/save", s.HandleRegister).Methods(POST)
	s.mux.HandleFunc("/api/products", s.HandleGetProduct).Methods(GET)
	s.mux.HandleFunc("/api/products/category", s.HandleGetCategory).Methods(GET)
	s.mux.HandleFunc("/api/products/{category}", s.HandleGetProductByCategory).Methods(GET)
	s.mux.HandleFunc("/api/products/search/{Word}", s.HandleGetProductBySearch).Methods(GET)
	s.mux.HandleFunc("/api/customer/login", s.HandleLogin).Methods(POST)
	SubRoutineCustomer.HandleFunc("/check", s.HandleChecker).Methods(POST)
	SubRoutineCustomer.HandleFunc("/delete", s.HandleDelete).Methods(DELETE)
	s.mux.HandleFunc("/api/customer/get/{ID}", s.HandleGetByID).Methods(GET)
	SubRoutineCustomer.HandleFunc("/save/products", s.HandleAddToCart).Methods(POST)

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

func (s *Server) HandleChecker(writer http.ResponseWriter, request *http.Request) {
	log.Println(*<-channel)
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
	err = s.cartsSvc.SaveToCart(request.Context(), id, Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
