package router

import (
	"Storage/internal/database"
	"Storage/internal/model"
	"encoding/json"
	"fmt"
	"strconv"

	"time"

	// "log/slog"
	"net/http"
)

type ProductHandler struct {
	db database.PostgresStorage
	// log *slog.Logger
}

func NewProducHandler(db database.PostgresStorage) *ProductHandler {
	return &ProductHandler{db: db}
}

func (db *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 5
	}
	if offset < 0 {
		offset = 0
	}
	product, err := db.db.GetProduct(limit, offset)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func (db *ProductHandler) GetProductsById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	product, err := db.db.GetProductById(id)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(product)
}

func (db *ProductHandler) DeleteProducts(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	isExist, err := db.db.CheckIdOnExist(id)
	if err != nil {
		fmt.Fprint(w, "error", err)
	}

	if !isExist {
		http.Error(w, "Product id not found", http.StatusInternalServerError)
	} else {
		db.db.DeleteProduct(id)
		fmt.Fprint(w, "Product ", id, " was be deleted")
	}

}

func (db *ProductHandler) CreateProducts(w http.ResponseWriter, r *http.Request) {
	var req model.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Rubles:      req.Rubles,
		Pennies:     req.Pennies,
		Quantity:    req.Quantity,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := db.db.CreateProduct(&product); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func (db *ProductHandler) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	isExist, err := db.db.CheckIdOnExist(id)
	if err != nil {
		fmt.Fprint(w, "error", err)
	}

	if !isExist {
		http.Error(w, "Product id not found", http.StatusInternalServerError)
	}
	var req model.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.ID, _ = strconv.Atoi(id)
	req.UpdatedAt = time.Now().UTC()

	if err := db.db.UpdateProduct(&req); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
