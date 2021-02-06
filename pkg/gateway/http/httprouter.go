package http

import (
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/controller"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/swaggo/http-swagger"
	"net/http"
	"time"
)

// NewHTTPRouterServer creates a new http router server
func NewHTTPRouterServer(listenAddr string, accCtrl controller.AccountController, authCtrl controller.AuthController, trfCtrl controller.TransferController, authUC usecase.AuthUseCase, idpRepo repository.IdempotencyRepository) *http.Server {
	router := httprouter.New()
	router.PanicHandler = handlePanic
	router.GlobalOPTIONS = http.HandlerFunc(handleOPTIONS)

	// accounts
	router.HandlerFunc(http.MethodPost, "/accounts", middleware.Idempotency(idpRepo, accCtrl.Create))
	router.HandlerFunc(http.MethodGet, "/accounts", accCtrl.Fetch)
	router.HandlerFunc(http.MethodGet, "/accounts/:id/balance", accCtrl.GetBalance)

	// auth
	router.HandlerFunc(http.MethodPost, "/login", authCtrl.Login)

	// transfer
	router.HandlerFunc(http.MethodPost, "/transfers", middleware.BearerAuth(authUC, middleware.Idempotency(idpRepo, trfCtrl.Create)))
	router.HandlerFunc(http.MethodGet, "/transfers", middleware.BearerAuth(authUC, trfCtrl.Fetch))

	router.HandlerFunc(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)
	router.HandlerFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger", http.StatusFound)
	})

	c := alice.New()
	c = c.Append(middleware.NewLoggerHandlerFunc())

	return &http.Server{
		Addr:         listenAddr,
		Handler:      c.Then(router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
