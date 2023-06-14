package handlers

import (
	"encoding/json"
	"github.com/RakhimovAns/Shop/types"
	"net/http"
)

func (s *Server) HandleAddToCart(writer http.ResponseWriter, request *http.Request) {
	var Products *[]types.Product
	err := json.NewDecoder(request.Body).Decode(&Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id := *<-channel
	err, ID := s.server.CartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.server.CartsSvc.SaveToCart(request.Context(), ID, Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = writer.Write([]byte("Added to cart successfully"))
}

// HandleDeleteProduct TODO Delete Product From Cart
func (s *Server) HandleDeleteProduct(writer http.ResponseWriter, request *http.Request) {
	var Products *[]types.Product
	err := json.NewDecoder(request.Body).Decode(&Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id := *<-channel
	err, ID := s.server.CartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.server.CartsSvc.DeleteProducts(request.Context(), ID, Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

// HandleChangeQTY TODO Change QTY in Cart
func (s *Server) HandleChangeQTY(writer http.ResponseWriter, request *http.Request) {
	var Products *[]types.Product
	err := json.NewDecoder(request.Body).Decode(&Products)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id := *<-channel
	err, ID := s.server.CartsSvc.GetCartID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.server.CartsSvc.ChangeQTY(request.Context(), ID, Products)
	if err == types.ErrNoSuch {
		http.Error(writer, "You want to delete more than you have ", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

// HandleGetCart TODO Get Cart
func (s *Server) HandleGetCart(writer http.ResponseWriter, request *http.Request) {
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
	data, err := json.Marshal(Products)
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

// HandleGetSum TODO Get Sum Of Cart
func (s *Server) HandleGetSum(writer http.ResponseWriter, request *http.Request) {
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
	data, err := json.Marshal(sum)
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

// HandleDeleteCart TODO Delete Cart
func (s *Server) HandleDeleteCart(writer http.ResponseWriter, request *http.Request) {
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
}
