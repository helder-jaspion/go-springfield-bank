package main

import (
	"github.com/helder-jaspion/go-springfield-bank/config"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/datasource/postgres"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/datasource/redis"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/controller"
	"github.com/helder-jaspion/go-springfield-bank/pkg/infraestructure/logging"
	"github.com/helder-jaspion/go-springfield-bank/pkg/infraestructure/monitoring"
	"github.com/rs/zerolog/log"
)

func main() {
	conf := config.ReadConfig("config/.env")

	logging.InitZeroLog(conf.Log.Level, conf.Log.Encoding)

	dbPool, err := postgres.ConnectPool(conf.Postgres.GetDSN(), conf.Postgres.Migrate)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("error connecting to db")
	}
	defer dbPool.Close()

	redisClient := redis.Connect(conf.Redis.Addr, conf.Redis.Password)
	defer func() {
		err = redisClient.Close()
		log.Fatal().Stack().Err(err).Msg("error closing redis connection")
	}()

	go monitoring.RunServer(conf.Monitoring.Port, dbPool, redisClient)

	accRepo := postgres.NewAccountRepository(dbPool)
	accUC := usecase.NewAccountUseCase(accRepo)
	accCtrl := controller.NewAccountController(accUC)

	authUC := usecase.NewAuthUseCase(conf.Auth.SecretKey, conf.Auth.AccessTokenDur, accRepo)
	authCtrl := controller.NewAuthController(authUC)

	trfRepo := postgres.NewTransferRepository(dbPool)
	trfUC := usecase.NewTransferUseCase(trfRepo, accRepo)
	trfCtrl := controller.NewTransferController(trfUC, authUC)

	idpRepo := redis.NewIdempotencyRepository(redisClient)

	httpRouterSrv := http.NewHTTPRouterServer(":"+conf.API.HTTPPort, accCtrl, authCtrl, trfCtrl, authUC, idpRepo)
	http.StartServer(httpRouterSrv)
}
