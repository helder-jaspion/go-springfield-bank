package main

import (
	"github.com/helder-jaspion/go-springfield-bank/pkg/adapter/controller"
	"github.com/helder-jaspion/go-springfield-bank/pkg/adapter/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	accountMemRepo := repository.NewAccountMemoryRepository()
	accountUC := usecase.NewAccountUseCase(accountMemRepo)
	accountController := controller.NewAccountController(accountUC)

	r := httprouter.New()
	r.HandlerFunc("POST", "/accounts", accountController.Create)

	log.Fatal(http.ListenAndServe(":8080", r))
}
