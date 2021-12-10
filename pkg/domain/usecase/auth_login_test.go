package usecase

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository/mock"
)

func Test_authUseCase_Login(t *testing.T) {
	t.Parallel()

	backgroundCtx := context.Background()

	type fields struct {
		secretKey      string
		accessTokenDur time.Duration
		accRepo        repository.AccountRepository
	}
	type args struct {
		ctx        context.Context
		loginInput AuthLoginInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantSub string
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				secretKey:      "whatever",
				accessTokenDur: 1 * time.Minute,
				accRepo: mock.AccountRepository{
					OnGetByCPF: func(ctx context.Context, cpf model.CPF) (*model.Account, error) {
						hashedSecret, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
						return &model.Account{
							ID:        "any-uuid-1",
							Name:      "Jon Snow",
							CPF:       "59951332099",
							Secret:    string(hashedSecret),
							Balance:   0,
							CreatedAt: time.Time{},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				loginInput: AuthLoginInput{
					CPF:    "59951332099",
					Secret: "secret",
				},
			},
			wantSub: "any-uuid-1",
			wantErr: nil,
		},
		{
			name: "wrong password should return invalid credentials error",
			fields: fields{
				secretKey:      "whatever",
				accessTokenDur: 1 * time.Minute,
				accRepo: mock.AccountRepository{
					OnGetByCPF: func(ctx context.Context, cpf model.CPF) (*model.Account, error) {
						hashedSecret, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
						return &model.Account{
							ID:        "any-uuid-1",
							Name:      "Jon Snow",
							CPF:       "59951332099",
							Secret:    string(hashedSecret),
							Balance:   0,
							CreatedAt: time.Time{},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				loginInput: AuthLoginInput{
					CPF:    "59951332099",
					Secret: "wrong",
				},
			},
			wantSub: "any-uuid-1",
			wantErr: ErrAuthInvalidCredentials,
		},
		{
			name: "not found account should return invalid credentials error",
			fields: fields{
				secretKey:      "whatever",
				accessTokenDur: 1 * time.Minute,
				accRepo: mock.AccountRepository{
					OnGetByCPF: func(ctx context.Context, cpf model.CPF) (*model.Account, error) {
						return nil, repository.ErrAccountNotFound
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				loginInput: AuthLoginInput{
					CPF:    "59951332099",
					Secret: "secret",
				},
			},
			wantSub: "any-uuid-1",
			wantErr: ErrAuthInvalidCredentials,
		},
		{
			name: "getByCPF returns no ID should return invalid credentials error",
			fields: fields{
				secretKey:      "whatever",
				accessTokenDur: 1 * time.Minute,
				accRepo: mock.AccountRepository{
					OnGetByCPF: func(ctx context.Context, cpf model.CPF) (*model.Account, error) {
						hashedSecret, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
						return &model.Account{
							Name:      "Jon Snow",
							CPF:       "59951332099",
							Secret:    string(hashedSecret),
							Balance:   0,
							CreatedAt: time.Time{},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				loginInput: AuthLoginInput{
					CPF:    "59951332099",
					Secret: "secret",
				},
			},
			wantSub: "any-uuid-1",
			wantErr: ErrAuthInvalidCredentials,
		},
		{
			name: "getByCPF other error should return login error",
			fields: fields{
				secretKey:      "whatever",
				accessTokenDur: 1 * time.Minute,
				accRepo: mock.AccountRepository{
					OnGetByCPF: func(ctx context.Context, cpf model.CPF) (*model.Account, error) {
						return nil, errors.New("any error")
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				loginInput: AuthLoginInput{
					CPF:    "59951332099",
					Secret: "secret",
				},
			},
			wantSub: "any-uuid-1",
			wantErr: ErrAuthLogin,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authUC := NewAuthUseCase(
				tt.fields.secretKey,
				tt.fields.accessTokenDur,
				tt.fields.accRepo,
			)

			got, err := authUC.Login(tt.args.ctx, tt.args.loginInput)
			if err != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			token, err := jwt.ParseWithClaims(
				got.AccessToken,
				&jwt.StandardClaims{},
				func(token *jwt.Token) (interface{}, error) {
					_, ok := token.Method.(*jwt.SigningMethodHMAC)
					if !ok {
						return nil, fmt.Errorf("unexpected token signing method")
					}

					return []byte(tt.fields.secretKey), nil
				},
			)
			if err != nil {
				t.Errorf("Login() invalid token error = %v", err)
				return
			}

			claims, ok := token.Claims.(*jwt.StandardClaims)
			if !ok {
				t.Errorf("Login() invalid token error = %v", err)
				return
			}

			if !reflect.DeepEqual(claims.Subject, tt.wantSub) {
				t.Errorf("Login() got = %v, want %v", claims.Subject, tt.wantSub)
			}
		})
	}
}
