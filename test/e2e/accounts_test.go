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

	"github.com/helder-jaspion/go-springfield-bank/config"
	httpGateway "github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http"
)

func Test_accounts_Fetch(t *testing.T) {
	ja := jsonassert.New(t)

	type fields struct {
		dbPool      *pgxpool.Pool
		redisClient *redis.Client
		authConf    config.ConfAuth
	}
	type args struct {
		path string
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
			name: "empty database should return empty json",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/accounts",
			},
			wantStatus: 200,
			want:       "[]",
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "one result success",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/accounts",
			},
			wantStatus: 200,
			want:       `[{"id": "<<PRESENCE>>", "name": "Bart Simpson", "cpf": "123.456.789-11", "balance": 5.96, "created_at": "<<PRESENCE>>"}]`,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
					uuid.NewString(),
					"Bart Simpson",
					"12345678911",
					"secret",
					596)
				if err != nil {
					t.Errorf("GET %s, error on runBefore = %v", args.path, err)
				}
			},
		},
		{
			name: "two results sort order success",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/accounts",
			},
			wantStatus: 200,
			want: `[
					{"id": "<<PRESENCE>>", "name": "Bart Simpson", "cpf": "123.456.789-11", "balance": 5.96, "created_at": "<<PRESENCE>>"},
					{"id": "<<PRESENCE>>", "name": "Homer Simpson", "cpf": "123.456.789-12", "balance": 1234.5, "created_at": "<<PRESENCE>>"}
				]`,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
					uuid.NewString(),
					"Bart Simpson",
					"12345678911",
					"secret",
					596)
				if err != nil {
					t.Errorf("GET %s, error on runBefore = %v", args.path, err)
				}

				_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
					uuid.NewString(),
					"Homer Simpson",
					"12345678912",
					"secret",
					123450)
				if err != nil {
					t.Errorf("GET %s, error on runBefore = %v", args.path, err)
				}
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

			res, err := http.Get(ts.URL + tt.args.path)
			if err != nil {
				t.Fatal(err)
			}

			// Check the response status code
			if statusCode := res.StatusCode; statusCode != tt.wantStatus {
				t.Errorf("GET %s, statusCode = %v, wantStatus %v", tt.args.path, statusCode, tt.wantStatus)
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

func Test_accounts_Create(t *testing.T) {
	ja := jsonassert.New(t)

	type fields struct {
		dbPool      *pgxpool.Pool
		redisClient *redis.Client
		authConf    config.ConfAuth
	}
	type args struct {
		path string
		body string
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
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/accounts",
				body: `{"name": " Bart Simpson ", "cpf": " 343.639.16206 ", "secret": "s3cr3t", "balance": 5678.96}`,
			},
			wantStatus: 201,
			want:       `{"id": "<<PRESENCE>>", "name": "Bart Simpson", "cpf": "343.639.162-06", "balance": 5678.96, "created_at": "<<PRESENCE>>"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "invalid cpf should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/accounts",
				body: `{"name": " Bart Simpson ", "cpf": "11122211122", "balance": 5678.96}`,
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "'cpf' is invalid"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "invalid name should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/accounts",
				body: `{"name": "A", "cpf": "343.639.162-06", "secret": "s3cr3t", "balance": 5678.96}`,
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "'name' must be between 2 and 100 characters in length"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "invalid secret should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/accounts",
				body: `{"name": " Bart Simpson ", "cpf": "343.639.162-06", "balance": 5678.96}`,
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "'secret' must be between 6 and 100 characters in length"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "negative balance should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/accounts",
				body: `{"name": "Bart Simpson", "cpf": "343.639.162-06", "secret": "s3cr3t", "balance": -1}`,
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "'balance' must be greater than or equal to zero"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "existing cpf should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/accounts",
				body: `{"name": "Bart Simpson", "cpf": "343.639.162-06", "secret": "s3cr3t", "balance": 1}`,
			},
			wantStatus: 409,
			want:       `{"code": 409, "message": "an account with this CPF already exists"}`,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
					uuid.NewString(),
					"Bart Simpson",
					"34363916206",
					"secret",
					596)
				if err != nil {
					t.Errorf("POST %s, error on runBefore = %v", args.path, err)
				}
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
				path: "/accounts",
				body: `{invalid}`,
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "error reading input"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		// error reading input
	}
	// {"code":400,"message":"'secret' must be between 6 and 100 characters in length"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			ts := httptest.NewServer(httpGateway.GetHTTPHandler(tt.fields.dbPool, tt.fields.redisClient, tt.fields.authConf))
			defer ts.Close()

			res, err := http.Post(ts.URL+tt.args.path, jsonContentType, strings.NewReader(tt.args.body))
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

func Test_accounts_Create_Idempotent(t *testing.T) {
	ja := jsonassert.New(t)

	type fields struct {
		dbPool      *pgxpool.Pool
		redisClient *redis.Client
		authConf    config.ConfAuth
	}
	type args struct {
		path   string
		header map[string][]string
		body   string
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
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/accounts",
				header: map[string][]string{
					"X-Idempotency-Key": {time.Now().String()},
				},
				body: `{"name": "Bart Simpson", "cpf": "343.639.162-06", "secret": "s3cr3t", "balance": 5678.96}`,
			},
			wantStatus: 201,
			want:       `{"id": "<<PRESENCE>>", "name": "Bart Simpson", "cpf": "343.639.162-06", "balance": 5678.96, "created_at": "<<PRESENCE>>"}`,
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

			testReq := func(check func(*http.Response)) {
				ts := httptest.NewServer(httpGateway.GetHTTPHandler(tt.fields.dbPool, tt.fields.redisClient, tt.fields.authConf))
				defer ts.Close()

				req, err := http.NewRequest(http.MethodPost, ts.URL+tt.args.path, strings.NewReader(tt.args.body))
				if err != nil {
					t.Fatal(err)
				}
				if tt.args.header != nil {
					req.Header = tt.args.header
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

				if check != nil {
					check(res)
				}
			}

			testReq(nil)

			testReq(func(res *http.Response) {
				idempotencyCacheHeader := res.Header.Get("X-Idempotency-Cache")
				if idempotencyCacheHeader != "HIT" {
					t.Errorf("POST %s, X-Idempotency-Cache = %v, want %v", tt.args.path, idempotencyCacheHeader, "HIT")
				}
			})
		})
	}
}

func Test_accounts_GetBalance(t *testing.T) {
	ja := jsonassert.New(t)

	type fields struct {
		dbPool      *pgxpool.Pool
		redisClient *redis.Client
		authConf    config.ConfAuth
	}
	type args struct {
		path func() string
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
			name: "empty database should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: func() string {
					return fmt.Sprintf("/accounts/%s/balance", uuid.NewString())
				},
			},
			wantStatus: 404,
			want:       `{"code": 404, "message": "account not found"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "success",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      "any-secret",
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: func() string {
					id := uuid.NewString()
					path := fmt.Sprintf("/accounts/%s/balance", id)

					_, err := testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id,
						"Bart Simpson",
						"34363916206",
						"secret",
						596)
					if err != nil {
						t.Errorf("POST %s, error on runBefore = %v", path, err)
					}

					return path
				},
			},
			wantStatus: 200,
			want:       `{"id": "<<PRESENCE>>", "balance": 5.96}`,
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

			path := tt.args.path()
			res, err := http.Get(ts.URL + path)
			if err != nil {
				t.Fatal(err)
			}

			// Check the response status code
			if statusCode := res.StatusCode; statusCode != tt.wantStatus {
				t.Errorf("GET %s, statusCode = %v, wantStatus %v", path, statusCode, tt.wantStatus)
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
