package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// HandleGetProduct  TODO Get All Products
func (s *Server) HandleGetProduct(writer http.ResponseWriter, request *http.Request) {
	items, err := s.server.ProductsSvc.AllActiveProducts(request.Context())
	if err != nil {
		http.Error(writer, "Can't Get All Products", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ") // Настройте отступы для форматирования JSON

	for i, item := range items {
		if i > 0 {
			_, err = writer.Write([]byte("\n")) // Добавляем символ новой строки перед каждым элементом (кроме первого)
		}

		if err := encoder.Encode(item); err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

// HandleGetCategory TODO Get All Categories
func (s *Server) HandleGetCategory(writer http.ResponseWriter, request *http.Request) {
	items, err := s.server.ProductsSvc.AllCategories(request.Context())
	if err != nil {
		http.Error(writer, "Can't Get All Categories", http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = writer.Write([]byte("There is all categories\n"))
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// HandleGetProductByCategory TODO Get Products By Category
func (s *Server) HandleGetProductByCategory(writer http.ResponseWriter, request *http.Request) {
	Category := mux.Vars(request)["category"]
	items, err := s.server.ProductsSvc.GetByCategory(request.Context(), Category)
	if err != nil {
		http.Error(writer, "Can't Get Product By Category", http.StatusBadRequest)
		return
	}
	if len(items) == 0 {
		_, err = writer.Write([]byte("There is  not any product in this category"))
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = writer.Write([]byte("There is all products in this category\n"))
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// HandleGetProductBySearch TODO Search Product
func (s *Server) HandleGetProductBySearch(writer http.ResponseWriter, request *http.Request) {
	Key := mux.Vars(request)["Word"]
	Key = "%" + Key + "%"
	items, err := s.server.ProductsSvc.Search(request.Context(), Key)
	if err != nil {
		http.Error(writer, "Can't Search ", http.StatusBadRequest)
	}
	if len(items) == 0 {
		_, err = writer.Write([]byte("There is  not any product for this key"))
		return
	}
	data, err := json.Marshal(items)
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
