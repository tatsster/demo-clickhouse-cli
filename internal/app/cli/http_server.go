package cli

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tikivn/clickhousectl/internal/pkg/api"
)

func NewHttpServer(
	tableHandler *api.TableHandler,
	columnHandler *api.ColumnHandler,
) *HttpServer {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Mount("/api/tables", tableHandler.Route())
	mux.Mount("/api/column", columnHandler.Route())

	server := &http.Server{
		Handler: mux,
	}
	return &HttpServer{server}
}

type HttpServer struct {
	*http.Server
}
