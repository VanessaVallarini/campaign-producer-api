package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VanessaVallarini/campaign-producer-api/internal/api"
	"github.com/VanessaVallarini/campaign-producer-api/internal/config"
	"github.com/VanessaVallarini/campaign-producer-api/internal/dao"
	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/VanessaVallarini/campaign-producer-api/internal/pkg/cache"
	"github.com/VanessaVallarini/campaign-producer-api/internal/pkg/kafka/producer"
	"github.com/VanessaVallarini/campaign-producer-api/internal/pkg/postgres"
	"github.com/VanessaVallarini/campaign-producer-api/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	easyzap "github.com/lockp111/go-easyzap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoTracer "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfig()

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	server.Pre(middleware.RemoveTrailingSlash())
	server.Use(echoTracer.Middleware())
	server.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	server.Use(middleware.Recover())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST},
	}))

	timeLocation, err := time.LoadLocation(cfg.TimeLocation)
	if err != nil {
		easyzap.Fatal(ctx, err, "failed to load timeLocation")
	}

	// dao
	pool := postgres.CreatePool(ctx, &cfg.Database)
	ownerDao := dao.NewOwnerDao(pool)

	// client
	localCache := cache.NewLocalMapService()
	kafkaProducer := producer.NewProducer(ctx, cfg.KafkaOwner, model.OwnerAvro)

	// service
	ownerService := service.NewOwnerService(ownerDao, localCache, kafkaProducer, timeLocation)

	// api
	api.NewOwner(ownerService).Register(server)

	// Start HTTP server
	go func() {
		easyzap.Info(ctx, "starting http worker server at "+cfg.ServerHost)
		err := server.Start(cfg.ServerHost)
		easyzap.Fatal(ctx, err, "failed to start server")
	}()

	meta := echo.New()
	meta.HideBanner = true
	meta.HidePort = true

	meta.GET("/prometheus", echo.WrapHandler(promhttp.Handler()))

	api.NewHealthCheck().Register(meta)

	// starts meta application
	go func() {
		easyzap.Info(ctx, "starting metadata worker server at "+cfg.MetaHost)
		err := meta.Start(cfg.MetaHost)
		easyzap.Fatal(ctx, err, "failed to start meta server")
	}()

	// listens for system signals to gracefully shutdown
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	switch <-signalChannel {
	case os.Interrupt:
		easyzap.Info(context.Background(), "received SIGINT, stopping...")
	case syscall.SIGTERM:
		easyzap.Info(context.Background(), "received SIGTERM, stopping...")
	}
}
