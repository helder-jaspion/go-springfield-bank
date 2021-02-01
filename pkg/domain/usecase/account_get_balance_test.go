package usecase

import (
	"context"
	"errors"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository/mock"
	"reflect"
	"testing"
	"time"
)

func Test_accountUseCase_GetBalance(t *testing.T) {
	t.Parallel()

	backgroundCtx := context.Background()

	type fields struct {
		accountRepo repository.AccountRepository
	}
	type args struct {
		ctx context.Context
		id  model.AccountID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AccountBalanceOutput
		wantErr error
	}{
		{
			name: "repo error should return error",
			fields: fields{
				accountRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						return nil, errors.New("any database error")
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				id:  "any-uuid-1",
			},
			want:    nil,
			wantErr: ErrAccountGetBalance,
		},
		{
			name: "repo not found error should return not found error",
			fields: fields{
				accountRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						return nil, repository.ErrAccountNotFound
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				id:  "any-uuid-1",
			},
			want:    nil,
			wantErr: repository.ErrAccountNotFound,
		},
		{
			name: "zero balance successful",
			fields: fields{
				accountRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						return &model.Account{
							ID:        "any-uuid-1",
							Name:      "Jon Snow",
							CPF:       "599.513.320-99",
							Secret:    "IAmNotSnow",
							Balance:   0,
							CreatedAt: time.Time{},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				id:  "any-uuid-1",
			},
			want: &AccountBalanceOutput{
				ID:      "any-uuid-1",
				Balance: 0,
			},
			wantErr: nil,
		},
		{
			name: "positive integer balance successful",
			fields: fields{
				accountRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						return &model.Account{
							ID:        "any-uuid-1",
							Name:      "Jon Snow",
							CPF:       "599.513.320-99",
							Secret:    "IAmNotSnow",
							Balance:   1000,
							CreatedAt: time.Time{},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				id:  "any-uuid-1",
			},
			want: &AccountBalanceOutput{
				ID:      "any-uuid-1",
				Balance: 10,
			},
			wantErr: nil,
		},
		{
			name: "negative integer balance successful",
			fields: fields{
				accountRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						return &model.Account{
							ID:        "any-uuid-1",
							Name:      "Jon Snow",
							CPF:       "599.513.320-99",
							Secret:    "IAmNotSnow",
							Balance:   19900,
							CreatedAt: time.Time{},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				id:  "any-uuid-1",
			},
			want: &AccountBalanceOutput{
				ID:      "any-uuid-1",
				Balance: 199,
			},
			wantErr: nil,
		},
		{
			name: "positive decimal balance successful",
			fields: fields{
				accountRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						return &model.Account{
							ID:        "any-uuid-1",
							Name:      "Jon Snow",
							CPF:       "599.513.320-99",
							Secret:    "IAmNotSnow",
							Balance:   55,
							CreatedAt: time.Time{},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				id:  "any-uuid-1",
			},
			want: &AccountBalanceOutput{
				ID:      "any-uuid-1",
				Balance: 0.55,
			},
			wantErr: nil,
		},
		{
			name: "negative decimal balance successful",
			fields: fields{
				accountRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						return &model.Account{
							ID:        "any-uuid-1",
							Name:      "Jon Snow",
							CPF:       "599.513.320-99",
							Secret:    "IAmNotSnow",
							Balance:   -155,
							CreatedAt: time.Time{},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				id:  "any-uuid-1",
			},
			want: &AccountBalanceOutput{
				ID:      "any-uuid-1",
				Balance: -1.55,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountUC := NewAccountUseCase(tt.fields.accountRepo)

			got, err := accountUC.GetBalance(tt.args.ctx, tt.args.id)
			if err != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}
