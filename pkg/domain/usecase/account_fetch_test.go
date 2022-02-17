package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository/mock"
)

func Test_accountUseCase_Fetch(t *testing.T) {
	t.Parallel()

	backgroundCtx := context.Background()

	type fields struct {
		accRepo repository.AccountRepository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []AccountFetchOutput
		wantErr bool
	}{
		{
			name: "repo fetch error should return error",
			fields: fields{
				accRepo: mock.AccountRepository{
					OnFetch: func(ctx context.Context) ([]model.Account, error) {
						return nil, errors.New("any database error")
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repo empty result should return empty result",
			fields: fields{
				accRepo: mock.AccountRepository{
					OnFetch: func(ctx context.Context) ([]model.Account, error) {
						return []model.Account{}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
			},
			want:    []AccountFetchOutput{},
			wantErr: false,
		},
		{
			name: "repo one result should return one result",
			fields: fields{
				accRepo: mock.AccountRepository{
					OnFetch: func(ctx context.Context) ([]model.Account, error) {
						return []model.Account{
							{
								ID:        "any-uuid-1",
								Name:      "Jon Snow",
								CPF:       "59951332099",
								Secret:    "IAmNotSnow",
								Balance:   0,
								CreatedAt: time.Time{},
							},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
			},
			want: []AccountFetchOutput{
				{AccountCreateOutput: AccountCreateOutput{
					ID:        "any-uuid-1",
					Name:      "Jon Snow",
					CPF:       "599.513.320-99",
					Balance:   0,
					CreatedAt: time.Time{},
				}},
			},
			wantErr: false,
		},
		{
			name: "repo two results should return two results",
			fields: fields{
				accRepo: mock.AccountRepository{
					OnFetch: func(ctx context.Context) ([]model.Account, error) {
						return []model.Account{
							{
								ID:        "any-uuid-1",
								Name:      "Homer Simpson",
								CPF:       "59951332099",
								Secret:    "Donuts",
								Balance:   -240,
								CreatedAt: time.Time{},
							},
							{
								ID:        "any-uuid-2",
								Name:      "Marge Simpson",
								CPF:       "84352262048",
								Secret:    "Blu3H4ir",
								Balance:   55,
								CreatedAt: time.Time{},
							},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
			},
			want: []AccountFetchOutput{
				{
					AccountCreateOutput: AccountCreateOutput{
						ID:        "any-uuid-1",
						Name:      "Homer Simpson",
						CPF:       "599.513.320-99",
						Balance:   -2.4,
						CreatedAt: time.Time{},
					},
				},
				{
					AccountCreateOutput: AccountCreateOutput{
						ID:        "any-uuid-2",
						Name:      "Marge Simpson",
						CPF:       "843.522.620-48",
						Balance:   0.55,
						CreatedAt: time.Time{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			accountUC := NewAccountUseCase(tt.fields.accRepo)

			got, err := accountUC.Fetch(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
