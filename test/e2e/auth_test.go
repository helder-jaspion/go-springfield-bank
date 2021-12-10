package e2e

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kinbiko/jsonassert"
	"golang.org/x/crypto/bcrypt"

	"github.com/helder-jaspion/go-springfield-bank/config"
	httpGateway "github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http"
)

func Test_Login(t *testing.T) {
	ja := jsonassert.New(t)

	authSecret := "any-secret"

	type fields struct {
		dbPool      *pgxpool.Pool
		redisClient *redis.Client
		authConf    config.ConfAuth
	}
	type args struct {
		path          string
		headerAndBody func() (map[string][]string, string)
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       string
		runBefore  func(args)
	}{
		{
			name: "success",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/login",
				headerAndBody: func() (map[string][]string, string) {
					cpf := "34363916206"
					secret := "super#s3cr3t!"
					hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
					if err != nil {
						t.Errorf("Login %s, error on building header = %v", "/transfers", err)
					}

					id1 := uuid.NewString()
					_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id1,
						"Bart Simpson",
						cpf,
						string(hashedSecret),
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return nil, fmt.Sprintf(`{"cpf":"%s", "secret": "%s"}`, cpf, secret)
				},
			},
			wantStatus: 200,
			want:       `{"access_token": "<<PRESENCE>>"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "wrong cpf should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/login",
				headerAndBody: func() (map[string][]string, string) {
					cpf := "34363916206"
					secret := "super#s3cr3t!"
					hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
					if err != nil {
						t.Errorf("Login %s, error on building header = %v", "/transfers", err)
					}

					id1 := uuid.NewString()
					_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id1,
						"Bart Simpson",
						cpf,
						string(hashedSecret),
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return nil, fmt.Sprintf(`{"cpf":"%s", "secret": "%s"}`, "12312312311", secret)
				},
			},
			wantStatus: 401,
			want:       `{"code":401,"message":"invalid credentials"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "wrong secret should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/login",
				headerAndBody: func() (map[string][]string, string) {
					cpf := "34363916206"
					secret := "super#s3cr3t!"
					hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
					if err != nil {
						t.Errorf("Login %s, error on building header = %v", "/transfers", err)
					}

					id1 := uuid.NewString()
					_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id1,
						"Bart Simpson",
						cpf,
						string(hashedSecret),
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return nil, fmt.Sprintf(`{"cpf":"%s", "secret": "%s"}`, cpf, "Not my secret")
				},
			},
			wantStatus: 401,
			want:       `{"code":401,"message":"invalid credentials"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "invalid input json should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/login",
				headerAndBody: func() (map[string][]string, string) {
					id1 := uuid.NewString()
					_, err := testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id1,
						"Bart Simpson",
						"34363916206",
						"s3cr3t",
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return nil, `{invalid}`
				},
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "error reading input"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			ts := httptest.NewServer(httpGateway.GetHTTPHandler(tt.fields.dbPool, tt.fields.redisClient, tt.fields.authConf))
			defer ts.Close()

			reqHeader, reqBody := tt.args.headerAndBody()

			req, err := http.NewRequest(http.MethodPost, ts.URL+tt.args.path, strings.NewReader(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			if reqHeader != nil {
				req.Header = reqHeader
			}
			req.Header.Add(contentType, jsonContentType)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			// Check the response status code
			if statusCode := res.StatusCode; statusCode != tt.wantStatus {
				t.Errorf("POST %s, statusCode = %v, wantStatus %v", tt.args.path, statusCode, tt.wantStatus)
			}

			body, err := ioutil.ReadAll(res.Body)
			_ = res.Body.Close()

			if err != nil {
				t.Fatal(err)
			}

			ja.Assertf(string(body), tt.want)
		})
	}
}
