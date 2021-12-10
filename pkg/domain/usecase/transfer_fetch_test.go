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

func Test_transferUseCase_Fetch(t *testing.T) {
	t.Parallel()

	backgroundCtx := context.Background()

	type fields struct {
		trfRepo repository.TransferRepository
		accRepo repository.AccountRepository
	}
	type args struct {
		ctx       context.Context
		accountID model.AccountID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []TransferFetchOutput
		wantErr bool
	}{
		{
			name: "repo fetch error should return error",
			fields: fields{
				trfRepo: mock.TransferRepository{
					OnFetch: func(ctx context.Context, accountID model.AccountID) ([]model.Transfer, error) {
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
				trfRepo: mock.TransferRepository{
					OnFetch: func(ctx context.Context, accountID model.AccountID) ([]model.Transfer, error) {
						return []model.Transfer{}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
			},
			want:    []TransferFetchOutput{},
			wantErr: false,
		},
		{
			name: "repo one result should return one result",
			fields: fields{
				trfRepo: mock.TransferRepository{
					OnFetch: func(ctx context.Context, accountID model.AccountID) ([]model.Transfer, error) {
						return []model.Transfer{
							{
								ID:                   "any-uuid-1",
								AccountOriginID:      "uuid-1",
								AccountDestinationID: "uuid-2",
								Amount:               1,
								CreatedAt:            time.Time{},
							},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
			},
			want: []TransferFetchOutput{
				{TransferCreateOutput: TransferCreateOutput{
					ID:                   "any-uuid-1",
					AccountOriginID:      "uuid-1",
					AccountDestinationID: "uuid-2",
					Amount:               0.01,
					CreatedAt:            time.Time{},
				}},
			},
			wantErr: false,
		},
		{
			name: "repo two results should return two results",
			fields: fields{
				trfRepo: mock.TransferRepository{
					OnFetch: func(ctx context.Context, accountID model.AccountID) ([]model.Transfer, error) {
						return []model.Transfer{
							{
								ID:                   "any-uuid-1",
								AccountOriginID:      "uuid-1",
								AccountDestinationID: "uuid-2",
								Amount:               1,
								CreatedAt:            time.Time{},
							},
							{
								ID:                   "any-uuid-2",
								AccountOriginID:      "uuid-3",
								AccountDestinationID: "uuid-4",
								Amount:               2,
								CreatedAt:            time.Time{},
							},
						}, nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
			},
			want: []TransferFetchOutput{
				{
					TransferCreateOutput: TransferCreateOutput{
						ID:                   "any-uuid-1",
						AccountOriginID:      "uuid-1",
						AccountDestinationID: "uuid-2",
						Amount:               0.01,
						CreatedAt:            time.Time{},
					},
				},
				{
					TransferCreateOutput: TransferCreateOutput{
						ID:                   "any-uuid-2",
						AccountOriginID:      "uuid-3",
						AccountDestinationID: "uuid-4",
						Amount:               0.02,
						CreatedAt:            time.Time{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trfUC := &transferUseCase{
				trfRepo: tt.fields.trfRepo,
				accRepo: tt.fields.accRepo,
			}
			got, err := trfUC.Fetch(tt.args.ctx, tt.args.accountID)
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
