package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"Storage/internal/model"
	"Storage/internal/storage"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrIDNotFound       = errors.New("id not found")
	ErrIncorrectID      = errors.New("incorrect id")
)

type ProductHandler struct {
	logger  *slog.Logger
	storage *storage.Storage
}

func NewProducHandler(logger *slog.Logger, storage *storage.Storage) *ProductHandler {
	return &ProductHandler{
		logger:  logger,
		storage: storage,
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Error("GetProduct (1): %w", "error", ErrMethodNotAllowed)
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)

		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 5
	}
	if offset < 0 {
		offset = 0
	}

	product, err := h.storage.GetProduct(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("GetProduct (2): %w", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	JSON(w, r, http.StatusOK, product)
}

func (h *ProductHandler) GetProductsById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("GetProductByID (1): %w", "error", err)
		http.Error(w, ErrIncorrectID.Error(), http.StatusBadRequest)

		return
	}
	product, err := h.storage.GetProductById(r.Context(), id)
	if err != nil {
		h.logger.Error("GetProductByID (2): %w", "error", err)
		http.Error(w, ErrIDNotFound.Error(), http.StatusInternalServerError)

		return
	}

	JSON(w, r, http.StatusOK, product)
}

func (h *ProductHandler) DeleteProducts(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("GetProductByID (1): %w", "error", err)
		http.Error(w, ErrIDNotFound.Error(), http.StatusBadRequest)
	}
	if err := h.storage.DeleteProduct(r.Context(), id); err != nil {
		h.logger.Error("DeleteProduct (2): %w", "error", err)
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)

		return
	}

	JSON(w, r, http.StatusOK, nil)
}

func (h *ProductHandler) CreateProducts(w http.ResponseWriter, r *http.Request) {
	var req model.Product

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("CreateProduct (1): %w", "error", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)

		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Error("CrateProduct (2): %w", "error", err)
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

	if err := h.storage.CreateProduct(r.Context(), &product); err != nil {
		h.logger.Error("CreateProduct (3): %w", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	JSON(w, r, http.StatusOK, product)
}

func (h *ProductHandler) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req model.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("UpdateProduct (1): %w", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	req.ID, _ = strconv.Atoi(id)
	req.UpdatedAt = time.Now().UTC()

	if err := h.storage.UpdateProduct(r.Context(), &req); err != nil {
		if errors.Is(err, ErrMethodNotAllowed) {
			h.logger.Error("UpdateProduct (2): %w", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

}
