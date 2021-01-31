package main

import (
	"github.com/helder-jaspion/go-springfield-bank/pkg/adapter/controller"
	"github.com/helder-jaspion/go-springfield-bank/pkg/adapter/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/infraestructure/delivery/http"
	"github.com/helder-jaspion/go-springfield-bank/pkg/infraestructure/logging"
)

func main() {
	logging.InitZerolog("debug", "json")

	accountMemRepo := repository.NewAccountMemoryRepository()
	accountUC := usecase.NewAccountUseCase(accountMemRepo)
	accountController := controller.NewAccountController(accountUC)

	httpRouterSrv := http.NewHTTPRouterServer(":8080", accountController)
	http.StartServer(httpRouterSrv)
}
