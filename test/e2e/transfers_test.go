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
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kinbiko/jsonassert"

	"github.com/helder-jaspion/go-springfield-bank/config"
	httpGateway "github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http"
)

func Test_transfers_Fetch(t *testing.T) {
	ja := jsonassert.New(t)

	authSecret := "any-secret"

	type fields struct {
		dbPool      *pgxpool.Pool
		redisClient *redis.Client
		authConf    config.ConfAuth
	}
	type args struct {
		path   string
		header func() map[string][]string
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
			name: "no Authorization header should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
				header: func() map[string][]string {
					return nil
				},
			},
			wantStatus: 401,
			want:       `{"code":401,"message":"malformed Token"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "invalid token should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
				header: func() map[string][]string {
					return map[string][]string{
						"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
					}
				},
			},
			wantStatus: 401,
			want:       `{"code":401,"message":"invalid access token"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "empty db should return empty",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
				header: func() map[string][]string {
					id := uuid.NewString()
					_, err := testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id,
						"Bart Simpson",
						"34363916206",
						"s3cr3t",
						596)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					now := time.Now()
					accessTokenClaims := jwt.RegisteredClaims{
						Subject:   id,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
					}

					accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
					accessTokenString, err := accessToken.SignedString([]byte(authSecret))
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return map[string][]string{
						"Authorization": {"Bearer " + accessTokenString},
					}
				},
			},
			wantStatus: 200,
			want:       `[]`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "should success",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
				header: func() map[string][]string {
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

					id2 := uuid.NewString()
					_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id2,
						"Homer Simpson",
						"62792172053",
						"s3cr3t2",
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					_, err = testDbPool.Exec(context.Background(), "INSERT INTO transfers (id, account_origin_id, account_destination_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
						uuid.NewString(),
						id1,
						id2,
						1,
						time.Now())
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					_, err = testDbPool.Exec(context.Background(), "INSERT INTO transfers (id, account_origin_id, account_destination_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
						uuid.NewString(),
						id2,
						id1,
						2,
						time.Now())
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					now := time.Now()
					accessTokenClaims := jwt.RegisteredClaims{
						Subject:   id1,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
					}

					accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
					accessTokenString, err := accessToken.SignedString([]byte(authSecret))
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return map[string][]string{
						"Authorization": {"Bearer " + accessTokenString},
					}
				},
			},
			wantStatus: 200,
			want: `[
					{"id": "<<PRESENCE>>", "account_origin_id":"<<PRESENCE>>", "account_destination_id":"<<PRESENCE>>", "amount": 0.02, "created_at": "<<PRESENCE>>"},
					{"id": "<<PRESENCE>>", "account_origin_id":"<<PRESENCE>>", "account_destination_id":"<<PRESENCE>>", "amount": 0.01, "created_at": "<<PRESENCE>>"}
				]`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			ts := httptest.NewServer(httpGateway.GetHTTPHandler(tt.fields.dbPool, tt.fields.redisClient, tt.fields.authConf))
			defer ts.Close()

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ts.URL+tt.args.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header = tt.args.header()

			res, err := http.DefaultClient.Do(req)
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

func Test_transfers_Create(t *testing.T) {
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
			name: "no Authorization header should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
				headerAndBody: func() (map[string][]string, string) {
					return nil, "{}"
				},
			},
			wantStatus: 401,
			want:       `{"code":401,"message":"malformed Token"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "invalid token should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
				headerAndBody: func() (map[string][]string, string) {
					return map[string][]string{
						"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
					}, "{}"
				},
			},
			wantStatus: 401,
			want:       `{"code":401,"message":"invalid access token"}`,
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
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
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

					id2 := uuid.NewString()
					_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id2,
						"Homer Simpson",
						"62792172053",
						"s3cr3t2",
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					now := time.Now()
					accessTokenClaims := jwt.RegisteredClaims{
						Subject:   id1,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
					}

					accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
					accessTokenString, err := accessToken.SignedString([]byte(authSecret))
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return map[string][]string{
						"Authorization": {"Bearer " + accessTokenString},
					}, fmt.Sprintf(`{"account_destination_id":"%s", "amount": %f}`, id2, 0.25)
				},
			},
			wantStatus: 201,
			want:       `{"id": "<<PRESENCE>>", "account_origin_id":"<<PRESENCE>>", "account_destination_id":"<<PRESENCE>>", "amount": 0.25, "created_at": "<<PRESENCE>>"}`,
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
				path: "/transfers",
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

					id2 := uuid.NewString()
					_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id2,
						"Homer Simpson",
						"62792172053",
						"s3cr3t2",
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					now := time.Now()
					accessTokenClaims := jwt.RegisteredClaims{
						Subject:   id1,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
					}

					accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
					accessTokenString, err := accessToken.SignedString([]byte(authSecret))
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return map[string][]string{
						"Authorization": {"Bearer " + accessTokenString},
					}, `{invalid}`
				},
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "error reading input"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "insufficient balance should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
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

					id2 := uuid.NewString()
					_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id2,
						"Homer Simpson",
						"62792172053",
						"s3cr3t2",
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					now := time.Now()
					accessTokenClaims := jwt.RegisteredClaims{
						Subject:   id1,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
					}

					accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
					accessTokenString, err := accessToken.SignedString([]byte(authSecret))
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return map[string][]string{
						"Authorization": {"Bearer " + accessTokenString},
					}, fmt.Sprintf(`{"account_destination_id":"%s", "amount": %f}`, id2, 100.01)
				},
			},
			wantStatus: 422,
			want:       `{"code": 422, "message": "current account balance is insufficient"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "negative amount should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
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

					id2 := uuid.NewString()
					_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id2,
						"Homer Simpson",
						"62792172053",
						"s3cr3t2",
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					now := time.Now()
					accessTokenClaims := jwt.RegisteredClaims{
						Subject:   id1,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
					}

					accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
					accessTokenString, err := accessToken.SignedString([]byte(authSecret))
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return map[string][]string{
						"Authorization": {"Bearer " + accessTokenString},
					}, fmt.Sprintf(`{"account_destination_id":"%s", "amount": %f}`, id2, -1.0)
				},
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "'amount' must be greater than zero"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "zero amount should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
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

					id2 := uuid.NewString()
					_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id2,
						"Homer Simpson",
						"62792172053",
						"s3cr3t2",
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					now := time.Now()
					accessTokenClaims := jwt.RegisteredClaims{
						Subject:   id1,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
					}

					accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
					accessTokenString, err := accessToken.SignedString([]byte(authSecret))
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return map[string][]string{
						"Authorization": {"Bearer " + accessTokenString},
					}, fmt.Sprintf(`{"account_destination_id":"%s", "amount": %f}`, id2, 0.0)
				},
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "'amount' must be greater than zero"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "non existing destination account should return error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
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

					id2 := uuid.NewString()

					now := time.Now()
					accessTokenClaims := jwt.RegisteredClaims{
						Subject:   id1,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
					}

					accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
					accessTokenString, err := accessToken.SignedString([]byte(authSecret))
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return map[string][]string{
						"Authorization": {"Bearer " + accessTokenString},
					}, fmt.Sprintf(`{"account_destination_id":"%s", "amount": %f}`, id2, 0.01)
				},
			},
			wantStatus: 422,
			want:       `{"code": 422, "message": "account not found"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "transfer to same account should error",
			fields: fields{
				dbPool:      testDbPool,
				redisClient: testRedisClient,
				authConf: config.ConfAuth{
					SecretKey:      authSecret,
					AccessTokenDur: 30 * time.Second,
				},
			},
			args: args{
				path: "/transfers",
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

					now := time.Now()
					accessTokenClaims := jwt.RegisteredClaims{
						Subject:   id1,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
					}

					accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
					accessTokenString, err := accessToken.SignedString([]byte(authSecret))
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return map[string][]string{
						"Authorization": {"Bearer " + accessTokenString},
					}, fmt.Sprintf(`{"account_destination_id":"%s", "amount": %f}`, id1, 0.25)
				},
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "origin and destination accounts must not be the same"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			ts := httptest.NewServer(httpGateway.GetHTTPHandler(tt.fields.dbPool, tt.fields.redisClient, tt.fields.authConf))
			defer ts.Close()

			reqHeader, reqBody := tt.args.headerAndBody()

			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, ts.URL+tt.args.path, strings.NewReader(reqBody))
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

func Test_transfers_Create_Idempotent(t *testing.T) {
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
				path: "/transfers",
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

					id2 := uuid.NewString()
					_, err = testDbPool.Exec(context.Background(), "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
						id2,
						"Homer Simpson",
						"62792172053",
						"s3cr3t2",
						10000)
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					now := time.Now()
					accessTokenClaims := jwt.RegisteredClaims{
						Subject:   id1,
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
					}

					accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
					accessTokenString, err := accessToken.SignedString([]byte(authSecret))
					if err != nil {
						t.Errorf("GET %s, error on building header = %v", "/transfers", err)
					}

					return map[string][]string{
						"Authorization":     {"Bearer " + accessTokenString},
						"X-Idempotency-Key": {time.Now().String()},
					}, fmt.Sprintf(`{"account_destination_id":"%s", "amount": %f}`, id2, 0.25)
				},
			},
			wantStatus: 201,
			want:       `{"id": "<<PRESENCE>>", "account_origin_id":"<<PRESENCE>>", "account_destination_id":"<<PRESENCE>>", "amount": 0.25, "created_at": "<<PRESENCE>>"}`,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			reqHeader, reqBody := tt.args.headerAndBody()

			testReq := func(check func(*http.Response)) {
				ts := httptest.NewServer(httpGateway.GetHTTPHandler(tt.fields.dbPool, tt.fields.redisClient, tt.fields.authConf))
				defer ts.Close()

				req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, ts.URL+tt.args.path, strings.NewReader(reqBody))
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
