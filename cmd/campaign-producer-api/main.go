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
	slugDao := dao.NewSlugDao(pool)
	slugHistoryDao := dao.NewSlugHistoryDao(pool)
	regionDao := dao.NewRegionDao(pool)
	regionHistoryDao := dao.NewRegionHistoryDao(pool)
	merchantDao := dao.NewMerchantDao(pool)
	campaignDao := dao.NewCampaignDao(pool)
	campaignHistoryDao := dao.NewCampaignHistoryDao(pool)
	spentDao := dao.NewSpentDao(pool)
	ledgerDao := dao.NewLedgerDao(pool)

	// client
	localCache := cache.NewLocalMapService()
	kafkaOwnerProducer := producer.NewProducer(ctx, cfg.KafkaOwner, model.OwnerAvro)
	kafkaSlugProducer := producer.NewProducer(ctx, cfg.KafkaSlug, model.SlugAvro)
	kafkaRegionProducer := producer.NewProducer(ctx, cfg.KafkaRegion, model.RegionAvro)
	kafkaMerchantProducer := producer.NewProducer(ctx, cfg.KafkaMerchant, model.MerchantAvro)
	kafkaCampaignProducer := producer.NewProducer(ctx, cfg.KafkaCampaign, model.CampaignAvro)
	kafkaSpentProducer := producer.NewProducer(ctx, cfg.KafkaSpent, model.SpentEventAvro)

	// service
	ownerService := service.NewOwnerService(ownerDao, localCache, kafkaOwnerProducer, timeLocation)
	slugService := service.NewSlugService(slugDao, slugHistoryDao, localCache, kafkaSlugProducer, timeLocation)
	regionService := service.NewRegionService(regionDao, regionHistoryDao, localCache, kafkaRegionProducer, timeLocation)
	merchantService := service.NewMerchantService(merchantDao, localCache, kafkaMerchantProducer, timeLocation)
	campaignService := service.NewCampaignService(campaignDao, campaignHistoryDao, localCache, kafkaCampaignProducer, timeLocation)
	spentService := service.NewSpentService(spentDao, kafkaSpentProducer, timeLocation)
	ledgerService := service.NewLedgerService(ledgerDao)

	// api
	api.NewOwner(ownerService).Register(server)
	api.NewSlug(slugService).Register(server)
	api.NewRegion(regionService).Register(server)
	api.NewMerchant(merchantService).Register(server)
	api.NewCampaign(campaignService).Register(server)
	api.NewSpent(spentService).Register(server)
	api.NewLedger(ledgerService).Register(server)

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
