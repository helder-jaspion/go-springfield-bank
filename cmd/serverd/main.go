package main

import (
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/db/memory"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/controller"
	"github.com/helder-jaspion/go-springfield-bank/pkg/infraestructure/logging"
)

func main() {
	logging.InitZerolog("debug", "json")

	accountMemRepo := memory.NewAccountRepository()
	accountUC := usecase.NewAccountUseCase(accountMemRepo)
	accountController := controller.NewAccountController(accountUC)

	httpRouterSrv := http.NewHTTPRouterServer(":8080", accountController)
	http.StartServer(httpRouterSrv)
}
