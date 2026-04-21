package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_middleware "github.com/zarhci/fulltodoap/internal/core/transport/http/middleware"
	core_server "github.com/zarhci/fulltodoap/internal/core/transport/http/server"
	users_transport_http "github.com/zarhci/fulltodoap/internal/feature/users/transport/http"
	"go.uber.org/zap"
)

func main() {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Error creating logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Starting application")

	users_transport_HTTP := users_transport_http.NewUsersHandler(nil)
	usersRouters := users_transport_HTTP.Router()

	apiVersionRouter := core_server.NewAPIVersionRouter(core_server.ApiVersionV1)
	apiVersionRouter.RegisterRoutes(usersRouters...)

	httpServer := core_server.NewServer(
		core_server.NewConfigMust(),
		logger,
		core_middleware.RequestID(),
		core_middleware.Logger(logger),
		core_middleware.PanicRecovery(),
		core_middleware.Trace(),
	)
	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Server error", zap.Error(err))
	}
}
