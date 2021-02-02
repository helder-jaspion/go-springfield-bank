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

	//accountRepo := memory.NewAccountRepository()
	accountRepo := postgres.NewAccountRepository(dbPool)
	accountUC := usecase.NewAccountUseCase(accountRepo)
	accountController := controller.NewAccountController(accountUC)

	httpRouterSrv := http.NewHTTPRouterServer(":"+conf.API.HTTPPort, accountController)
	http.StartServer(httpRouterSrv)
}
