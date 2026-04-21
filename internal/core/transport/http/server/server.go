package core_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_middleware "github.com/zarhci/fulltodoap/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type Server struct {
	mux        *http.ServeMux
	config     Config
	log        *core_logger.Logger
	middleware []core_middleware.Middleware
}

func NewServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_middleware.Middleware,
) *Server {
	return &Server{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (h *Server) RegisterAPIRoutes(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

func (h *Server) Run(ctx context.Context) error {
	mux := core_middleware.ChainMiddleware(h.mux, h.middleware...)

	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		h.log.Warn("starting server", zap.String("addr", h.config.Addr))
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server http: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), h.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown server: %w", err)
		}

		h.log.Warn("http server stopped")
	}

	return nil
}
