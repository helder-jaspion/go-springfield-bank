package main

import (
	"github.com/helder-jaspion/go-springfield-bank/config"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/db/memory"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/controller"
	"github.com/helder-jaspion/go-springfield-bank/pkg/infraestructure/logging"
)

func main() {
	conf := config.ReadConfigFromFile("config/.env")

	logging.InitZerolog(conf.Log.Level, conf.Log.Encoding)

	accountMemRepo := memory.NewAccountRepository()
	accountUC := usecase.NewAccountUseCase(accountMemRepo)
	accountController := controller.NewAccountController(accountUC)

	httpRouterSrv := http.NewHTTPRouterServer(":"+conf.API.HTTPPort, accountController)
	http.StartServer(httpRouterSrv)
}
