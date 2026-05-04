package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	core_postgres_pool "github.com/zarhci/fulltodoap/internal/core/repository/postgres/pool"
	core_middleware "github.com/zarhci/fulltodoap/internal/core/transport/http/middleware"
	core_server "github.com/zarhci/fulltodoap/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/zarhci/fulltodoap/internal/feature/tasks/repository/postgres"
	tasks_service "github.com/zarhci/fulltodoap/internal/feature/tasks/service"
	tasks_transport "github.com/zarhci/fulltodoap/internal/feature/tasks/transport/http"
	users_repository_postgres "github.com/zarhci/fulltodoap/internal/feature/users/repository/postgres"
	users_service "github.com/zarhci/fulltodoap/internal/feature/users/service"
	users_transport_http "github.com/zarhci/fulltodoap/internal/feature/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {

	time.Local = timeZone

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

	logger.Debug("application time zone", zap.Any("time_zone", timeZone))

	logger.Debug("Creating database connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("Error creating database connection pool:", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("Initializing feature", zap.String("feature", "users"))
	usersRepository := users_repository_postgres.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)

	logger.Debug("initializing feature", zap.String("feature", "tasks"), zap.String("repository", "postgres"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTaskService(tasksRepository)
	tasksTranport := tasks_transport.NewTasksHandler(tasksService)

	users_transport_HTTP := users_transport_http.NewUsersHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpServer := core_server.NewServer(
		core_server.NewConfigMust(),
		logger,
		core_middleware.RequestID(),
		core_middleware.Logger(logger),
		core_middleware.Trace(),
		core_middleware.PanicRecovery(),
	)

	apiVersionRouter := core_server.NewAPIVersionRouter(core_server.ApiVersionV1)
	apiVersionRouter.RegisterRoutes(users_transport_HTTP.Router()...)
	apiVersionRouter.RegisterRoutes(tasksTranport.Routes()...)
	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Server error", zap.Error(err))
	}
}
