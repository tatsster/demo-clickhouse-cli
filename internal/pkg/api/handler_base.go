package api

import (
	"net/http"

	"github.com/go-chi/render"
)

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

type BaseHandler struct {
}

func (h *BaseHandler) errorHandler(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	render.Status(r, statusCode)
	render.JSON(w, r, map[string]interface{}{
		"error":   err.Error(),
		"success": false,
		"data":    nil,
	})
}

func (h *BaseHandler) successHandler(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"error":   nil,
		"success": true,
		"data":    data,
	})
}
