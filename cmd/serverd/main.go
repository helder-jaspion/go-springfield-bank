package main

import (
	"github.com/helder-jaspion/go-springfield-bank/config"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/db/postgres"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/controller"
	"github.com/helder-jaspion/go-springfield-bank/pkg/infraestructure/logging"
)

func main() {
	conf := config.ReadConfigFromFile("config/.env")

	logging.InitZerolog(conf.Log.Level, conf.Log.Encoding)

	dbPool := postgres.ConnectPool(conf.Postgres.GetDSN(), conf.Postgres.Migrate)
	defer dbPool.Close()

	//accRepo := memory.NewAccountRepository()
	accRepo := postgres.NewAccountRepository(dbPool)
	accUC := usecase.NewAccountUseCase(accRepo)
	accCtrl := controller.NewAccountController(accUC)

	authUC := usecase.NewAuthUseCase(conf.Auth.SecretKey, conf.Auth.AccessTokenDur, accRepo)
	authCtrl := controller.NewAuthController(authUC)

	httpRouterSrv := http.NewHTTPRouterServer(":"+conf.API.HTTPPort, accCtrl, authCtrl)
	http.StartServer(httpRouterSrv)
}
