package main

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"syscall"

	"github.com/mbatimel/RegressionAnalysis/internal/config"
	"github.com/mbatimel/RegressionAnalysis/internal/service"
	transportHttp "github.com/mbatimel/RegressionAnalysis/internal/transport/http"
	"github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/externalapi"
	"github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/middlewares"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

const serviceName = "regression"

func main() {
	log.Logger = config.Values().Logger().With().Str("serviceName", serviceName).Logger()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	// redisClient := redis.NewClusterClient(&redis.ClusterOptions{
	// 		Addrs:      config.Values().RedisAddrs,
	// 		Password:   config.Values().RedisPassword,
	// 		MaxRetries: config.Values().RedisMaxRetries,
	// 		ReadOnly:   false,
	// 	})
	// 	redisStorage, err := redisInternal.New(redisClient)
	// 	if err != nil {
	// 			log.Logger.Fatal().Err(err).Msg("failed to connect to redis")
	// 		}

	// locker, err := locker.NewLocker(redisClient)
	// if err != nil {
	// 	log.Logger.Fatal().Err(err).Msg("failed to create locker")
	// }
	// postgresStorage, err := postgres.New(config.Values().Postgres, log.Logger)
	// if err != nil {
	// 	log.Logger.Fatal().Err(err).Msg("failed to connect to postgres")
	// }
	svc := service.Newservice(log.Logger)

	// innerServiceIDs := make(map[uuid.UUID]struct{}, len(config.Values().InnerServiceIDs))
	// for _, id := range config.Values().InnerServiceIDs {
	// 	innerServiceIDs[id] = struct{}{}
	// }

	// svc = middlewares.NewInternalMiddleware(svc, innerServiceIDs)

	services := []externalapi.Option{
		externalapi.Use(middlewares.Recover),
		externalapi.Regression(externalapi.NewRegression(svc)),
	}
	app := externalapi.New(log.Logger, services...).WithLog().WithMetrics()
	server := &fasthttp.Server{
		Handler:            app.Fiber().Handler(),
		MaxRequestBodySize: config.Values().MaxRequestBodySize,
		ReadBufferSize:     config.Values().MaxRequestHeaderSize,
		ReadTimeout:        time.Duration(config.Values().ReadTimeout) * time.Second,
	}

	healthServer := transportHttp.NewHealthServer()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.ServeMetrics(log.Logger, config.Values().MetricsPath, config.Values().MetricsBind)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		serveErr := server.ListenAndServe(config.Values().ServiceBind)
		if serveErr != nil {
			log.Fatal().Err(serveErr).Msg("failed to listen and serve pay-api-internal server")
		} else {
			log.Error().Msg("external api pay-api-internal server stopped with no error")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		healthErr := healthServer.Start(config.Values().HealthBind)
		if healthErr != nil {
			log.Error().Err(healthErr).Msg("failed to start health server")
		} else {
			log.Error().Msg("health server stopped with no error")
		}
	}()

	<-shutdown
	err := healthServer.Stop()
	if err != nil {
		log.Error().Err(err).Msg("failed to stop health server")
	}

	err = server.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("failed to shutdown server")
	}

	wg.Wait()
}
