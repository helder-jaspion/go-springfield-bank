package http

import (
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/controller"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"time"
)

// NewHTTPRouterServer creates a new http router server
func NewHTTPRouterServer(listenAddr string, accountController controller.AccountController) *http.Server {
	router := httprouter.New()
	router.PanicHandler = handlePanic
	router.GlobalOPTIONS = http.HandlerFunc(handleOPTIONS)

	router.HandlerFunc("POST", "/accounts", accountController.Create)
	router.HandlerFunc("GET", "/accounts", accountController.Fetch)
	router.HandlerFunc("GET", "/accounts/:id/balance", accountController.GetBalance)

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
