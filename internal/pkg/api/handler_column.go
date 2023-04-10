package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tikivn/clickhousectl/internal/pkg/entity/operator"
	"github.com/tikivn/clickhousectl/internal/pkg/service"
)

func NewColumnHandler(svc service.ColumnOperatorService) *ColumnHandler {
	return &ColumnHandler{svc: svc}
}

type ColumnHandler struct {
	*BaseHandler
	svc service.ColumnOperatorService
}

func (h *ColumnHandler) Route() chi.Router {
	router := chi.NewRouter()

	router.Post("/add", h.handleAddColumn)
	router.Delete("/drop", h.handleDropColumn)

	return router
}

func (h *ColumnHandler) handleAddColumn(w http.ResponseWriter, r *http.Request) {
	var req operator.OperatorAddColumn

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, err)
		return
	} else if err := req.Validate(); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Execute(r.Context(), &req); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(w, r, nil)
}

func (h *ColumnHandler) handleDropColumn(w http.ResponseWriter, r *http.Request) {
	var req operator.OperatorDropColumn

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, err)
		return
	} else if err := req.Validate(); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Execute(r.Context(), &req); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(w, r, nil)
}
