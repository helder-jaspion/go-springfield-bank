// @title GO Springfield Bank API
// @version 0.0.1
// @description GO Springfield Bank API simulates a digital bank where you can create and fetch accounts, login with your account and transfer money to other accounts.
// @description ### Authorization
// @description You can get the access_token returned from `/login`, click the **Authorize** button and input this format `Bearer <access_token>`. After this, the `Authorization` header will be sent along in your next requests.
// @description The JWT access token has short expiration, so maybe you have to log in again to get a new `access_token`.
// @description ### X-Idempotency-Key
// @description If you send the `X-Idempotency-Key` header along with a request, that request's response will be cached. So, if you send the same request with the same `X-Idempotency-Key` again, the server will respond the cached response, so no processing will be done twice.

// @contact.name Helder Alves
// @contact.url https://github.com/helder-jaspion/go-springfield-bank/

// @license.name MIT
// @license.url https://github.com/helder-jaspion/go-springfield-bank/blob/main/LICENSE

// @securityDefinitions.apikey Access token
// @in header
// @name Authorization

package main

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/helder-jaspion/go-springfield-bank/api"
	"github.com/helder-jaspion/go-springfield-bank/config"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/datasource/postgres"
	redisGateway "github.com/helder-jaspion/go-springfield-bank/pkg/gateway/datasource/redis"
	httpGateway "github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http"
	"github.com/helder-jaspion/go-springfield-bank/pkg/infraestructure/logging"
	"github.com/helder-jaspion/go-springfield-bank/pkg/infraestructure/monitoring"
)

func main() {
	conf := config.ReadConfig("config/.env")

	logging.InitZeroLog(conf.Log.Level, conf.Log.Encoding)

	dbPool, err := postgres.ConnectPool(conf.Postgres)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("error connecting to db")
	}
	defer dbPool.Close()

	redisClient, err := redisGateway.Connect(conf.Redis.URL)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("error connecting to redis")
	}
	defer func() {
		err = redisClient.Close()
		log.Fatal().Stack().Err(err).Msg("error closing redis connection")
	}()

	go monitoring.RunServer(conf.Monitoring.Port, dbPool, redisClient)

	api.SwaggerInfo.Host = conf.API.Host

	handler := httpGateway.GetHTTPHandler(dbPool, redisClient, conf.Auth)
	server := &http.Server{
		Addr:         ":" + conf.API.Port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	httpGateway.StartServer(server)
}
