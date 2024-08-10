package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/desepticon55/gofemart/internal/api/auth"
	customMiddleware "github.com/desepticon55/gofemart/internal/api/middleware"
	"github.com/desepticon55/gofemart/internal/api/order"
	"github.com/desepticon55/gofemart/internal/common"
	ordSrv "github.com/desepticon55/gofemart/internal/service/order"
	usrSrv "github.com/desepticon55/gofemart/internal/service/user"
	"github.com/desepticon55/gofemart/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

func main() {
	logger := initLogger()
	defer logger.Sync()

	config := parseConfig()

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(customMiddleware.CompressingMiddleware())
	router.Use(customMiddleware.DecompressingMiddleware())

	pool, err := createConnectionPool(context.Background(), config.DatabaseConnString)
	if err != nil {
		panic(err)
	}
	runMigrations(config.DatabaseConnString, logger)

	userRepository := storage.NewUserRepository(pool, logger)
	userService := usrSrv.NewUserService(logger, userRepository)

	orderRepository := storage.NewOrderRepository(pool, logger)
	orderService := ordSrv.NewOrderService(logger, orderRepository)

	router.Method(http.MethodPost, "/api/user/register", auth.RegisterHandler(logger, userService)) //регистрация пользователя
	router.Method(http.MethodPost, "/api/user/login", auth.LoginHandler(logger, userService))       //аутентификация пользователя

	router.Group(func(r chi.Router) {
		r.Use(customMiddleware.AuthMiddleware(logger))
		r.Method(http.MethodPost, "/api/user/orders", order.UploadOrderHandler(logger, orderService)) //загрузка пользователем номера заказа для расчёта
		//r.Method(http.MethodPost, "/api/user/balance/withdraw", metricsApi.NewFindAllMetricsHandler(metricsService, logger)) //запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
		r.Method(http.MethodGet, "/api/user/orders", order.FindAllOrdersHandler(logger, orderService)) //получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
		//r.Method(http.MethodGet, "/api/user/balance", metricsApi.NewFindAllMetricsHandler(metricsService, logger))           //получение текущего баланса счёта баллов лояльности пользователя
		//r.Method(http.MethodGet, "/api/user/balance/withdraw", metricsApi.NewFindAllMetricsHandler(metricsService, logger))  //запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
		//r.Method(http.MethodGet, "/api/user/withdrawals", metricsApi.NewFindAllMetricsHandler(metricsService, logger))       //получение информации о выводе средств с накопительного счёта пользователем
	})

	http.ListenAndServe(config.ServerAddress, router)
}

func createConnectionPool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("error parsing database config: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return pool, nil
}

func runMigrations(connectionString string, logger *zap.Logger) {
	databaseConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		logger.Error("Error during parse database URL", zap.Error(err))
		return
	}
	db := stdlib.OpenDB(*databaseConfig)
	defer db.Close()

	goose.SetDialect("postgres")
	if err := goose.Up(db, "migrations"); err != nil {
		logger.Error("Error during run database migrations", zap.Error(err))
	}
}

func parseConfig() common.Config {
	config := common.ParseConfig()
	flag.Parse()
	return config
}

func initLogger() *zap.Logger {
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)
	productionConfig := zap.NewProductionConfig()
	productionConfig.Encoding = "console"
	productionConfig.Level = level
	productionConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := productionConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
