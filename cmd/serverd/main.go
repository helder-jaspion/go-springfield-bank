package main

import (
	"github.com/helder-jaspion/go-springfield-bank/config"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/db/postgres"
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
		log.Fatal().Err(err).Msg("error connecting to db")
	}
	defer dbPool.Close()

	go monitoring.RunServer(conf.Monitoring.Port, dbPool)

	accRepo := postgres.NewAccountRepository(dbPool)
	accUC := usecase.NewAccountUseCase(accRepo)
	accCtrl := controller.NewAccountController(accUC)

	authUC := usecase.NewAuthUseCase(conf.Auth.SecretKey, conf.Auth.AccessTokenDur, accRepo)
	authCtrl := controller.NewAuthController(authUC)

	trfRepo := postgres.NewTransferRepository(dbPool)
	trfUC := usecase.NewTransferUseCase(trfRepo, accRepo)
	trfCtrl := controller.NewTransferController(trfUC, authUC)

	httpRouterSrv := http.NewHTTPRouterServer(":"+conf.API.HTTPPort, accCtrl, authCtrl, trfCtrl, authUC)
	http.StartServer(httpRouterSrv)
}

// TODO log withStack
// TODO logout
// TODO grpc
