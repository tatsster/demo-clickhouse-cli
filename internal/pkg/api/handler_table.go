package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tikivn/clickhousectl/internal/pkg/entity/operator"
	"github.com/tikivn/clickhousectl/internal/pkg/service"
)

func NewTableHandler(svc service.TableOperatorService) *TableHandler {
	return &TableHandler{svc: svc}
}

type TableHandler struct {
	*BaseHandler
	svc service.TableOperatorService
}

func (h *TableHandler) Route() chi.Router {
	router := chi.NewRouter()

	router.Post("/create", h.handleCreateTable)
	router.Delete("/drop", h.handleDropTable)

	return router
}

func (h *TableHandler) handleCreateTable(w http.ResponseWriter, r *http.Request) {
	var req operator.OperatorCreateTable

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

func (h *TableHandler) handleDropTable(w http.ResponseWriter, r *http.Request) {
	var req operator.OperatorDropTable

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
