package http

import (
	"github.com/go-redis/redis"
	"github.com/helder-jaspion/go-springfield-bank/config"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/datasource/postgres"
	redisGateway "github.com/helder-jaspion/go-springfield-bank/pkg/gateway/datasource/redis"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/controller"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/swaggo/http-swagger"
	"net/http"
)

// NewHTTPRouterHandler creates a new http router handler.
func NewHTTPRouterHandler(
	accCtrl controller.AccountController,
	authCtrl controller.AuthController,
	trfCtrl controller.TransferController,
	authUC usecase.AuthUseCase,
	idpRepo repository.IdempotencyRepository,
) http.Handler {
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

	return c.Then(router)
}

// GetHTTPHandler instantiates the repos, ucs and controllers and returns a handler.
func GetHTTPHandler(dbPool *pgxpool.Pool, redisClient *redis.Client, authConf config.ConfAuth) http.Handler {
	accRepo := postgres.NewAccountRepository(dbPool)
	accUC := usecase.NewAccountUseCase(accRepo)
	accCtrl := controller.NewAccountController(accUC)

	authUC := usecase.NewAuthUseCase(authConf.SecretKey, authConf.AccessTokenDur, accRepo)
	authCtrl := controller.NewAuthController(authUC)

	trfRepo := postgres.NewTransferRepository(dbPool)
	trfUC := usecase.NewTransferUseCase(trfRepo, accRepo)
	trfCtrl := controller.NewTransferController(trfUC, authUC)

	idpRepo := redisGateway.NewIdempotencyRepository(redisClient)

	return NewHTTPRouterHandler(accCtrl, authCtrl, trfCtrl, authUC, idpRepo)
}
