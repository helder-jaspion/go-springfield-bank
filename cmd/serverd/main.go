// @title GO Springfield Bank API
// @version 0.0.1
// @description GO Springfield Bank API simulates a digital bank where you can create and fetch accounts, login with your account and transfer money to other accounts.

// @contact.name Helder Alves
// @contact.email helder.jaspion@gmail.com
// @contact.url https://github.com/helder-jaspion/go-springfield-bank/

// @license.name MIT
// @license.url https://github.com/helder-jaspion/go-springfield-bank/blob/main/LICENSE

// @securityDefinitions.apikey Access token
// @in header
// @name Authorization

package main

import (
	"github.com/helder-jaspion/go-springfield-bank/api"
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

	redisClient, err := redis.Connect(conf.Redis.URL)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("error connecting to redis")
	}
	defer func() {
		err = redisClient.Close()
		log.Fatal().Stack().Err(err).Msg("error closing redis connection")
	}()

	go monitoring.RunServer(conf.Monitoring.Port, dbPool, redisClient)

	api.SwaggerInfo.Host = conf.API.Host

	accRepo := postgres.NewAccountRepository(dbPool)
	accUC := usecase.NewAccountUseCase(accRepo)
	accCtrl := controller.NewAccountController(accUC)

	authUC := usecase.NewAuthUseCase(conf.Auth.SecretKey, conf.Auth.AccessTokenDur, accRepo)
	authCtrl := controller.NewAuthController(authUC)

	trfRepo := postgres.NewTransferRepository(dbPool)
	trfUC := usecase.NewTransferUseCase(trfRepo, accRepo)
	trfCtrl := controller.NewTransferController(trfUC, authUC)

	idpRepo := redis.NewIdempotencyRepository(redisClient)

	httpRouterSrv := http.NewHTTPRouterServer(":"+conf.API.Port, accCtrl, authCtrl, trfCtrl, authUC, idpRepo)
	http.StartServer(httpRouterSrv)
}
