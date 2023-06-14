package handlers

import (
	"encoding/json"
	"github.com/RakhimovAns/Shop/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// HandleRegister TODO  Register Customer
func (s *Server) HandleRegister(writer http.ResponseWriter, request *http.Request) {
	var item *types.Customer
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.server.CustomersSvc.Register(request.Context(), item)
	if err != nil {
		http.Error(writer, "Register has been failed", http.StatusInternalServerError)
		return
	}
	_, err = writer.Write([]byte("Was saved successfully"))
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleLogin(writer http.ResponseWriter, request *http.Request) {
	var item *types.Customer
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	token, err := s.server.CustomersSvc.Login(request.Context(), item.Phone, *item.Password)
	if err == types.ErrNoSuchUser {
		_, err = writer.Write([]byte("No account with this phone number"))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if err == types.ErrInvalidPassword {
		_, err = writer.Write([]byte("Passwords don't match"))
		if err != nil {
			http.Error(writer, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	} else if err == nil {
		_, err = writer.Write([]byte("You have been login successfully\nHere is your Token\n"))
		if err != nil {
			http.Error(writer, http.StatusText(500), 500)
			return
		}
		data, err := json.Marshal(token)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(data)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

// HandleMakePurchase TODO Make Purchase
func (s *Server) HandleMakePurchase(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	err, ID := s.server.CartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	Products, err := s.server.CartsSvc.GetCartBYID(request.Context(), ID)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	sum, err := s.server.CartsSvc.GetSum(request.Context(), Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err, Customer := s.server.CustomersSvc.GetByID(request.Context(), id)
	if Customer.Balance < sum {
		http.Error(writer, "You don't have enough money to buy it", http.StatusBadRequest)
		return
	}
	err = s.server.PurchasesSvc.AddToPurchase(request.Context(), Products, id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.server.CustomersSvc.ChangeBalance(request.Context(), id, sum)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.server.CartsSvc.DeleteCart(request.Context(), ID)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, _ = writer.Write([]byte("You have bought it"))
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
	err = s.server.CustomersSvc.DepositBalance(request.Context(), id, item.Balance)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

// HandleDelete TODO Delete Customer
func (s *Server) HandleDelete(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	err, ID := s.server.CartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.server.CartsSvc.DeleteCart(request.Context(), ID)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.server.PurchasesSvc.DeletePurchase(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.server.CustomersSvc.Delete(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = writer.Write([]byte("Was deleted"))
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
	err, Customer := s.server.CustomersSvc.GetByID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(Customer)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
