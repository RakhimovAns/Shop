package handlers

import (
	"encoding/json"
	"net/http"
)

// HandleGetAllPurchaseByID TODO Get All Purchases By ID
func (s *Server) HandleGetAllPurchaseByID(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	Purchases, err := s.server.PurchasesSvc.GetAllPurchase(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(Purchases)
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
