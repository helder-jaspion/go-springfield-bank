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

func TestTransferCreateInput_Validate(t *testing.T) {
	t.Parallel()

	type fields struct {
		AccountOriginID      string
		AccountDestinationID string
		Amount               float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "empty origin account should return error",
			fields: fields{
				AccountDestinationID: "uuid-2",
				Amount:               10,
			},
			wantErr: ErrTransferOriginAccountRequired,
		},
		{
			name: "empty destination account should return error",
			fields: fields{
				AccountOriginID: "uuid-1",
				Amount:          10,
			},
			wantErr: ErrTransferDestinationAccountRequired,
		},
		{
			name: "zero amount should return error",
			fields: fields{
				AccountOriginID:      "uuid-1",
				AccountDestinationID: "uuid-2",
				Amount:               0,
			},
			wantErr: ErrTransferAmountNotPositive,
		},
		{
			name: "negative amount should return error",
			fields: fields{
				AccountOriginID:      "uuid-1",
				AccountDestinationID: "uuid-2",
				Amount:               -1,
			},
			wantErr: ErrTransferAmountNotPositive,
		},
		{
			name: "same account should return error",
			fields: fields{
				AccountOriginID:      "uuid-1",
				AccountDestinationID: "uuid-1",
				Amount:               10,
			},
			wantErr: ErrTransferSameAccount,
		},
		{
			name: "success",
			fields: fields{
				AccountOriginID:      "uuid-1",
				AccountDestinationID: "uuid-2",
				Amount:               10,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &TransferCreateInput{
				AccountOriginID:      tt.fields.AccountOriginID,
				AccountDestinationID: tt.fields.AccountDestinationID,
				Amount:               tt.fields.Amount,
			}
			if err := input.Validate(); err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_transferUseCase_Create(t *testing.T) {
	t.Parallel()

	backgroundCtx := context.Background()

	type fields struct {
		trfRepo repository.TransferRepository
		accRepo repository.AccountRepository
	}
	type args struct {
		ctx           context.Context
		transferInput TransferCreateInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TransferCreateOutput
		wantErr error
	}{
		{
			name: "zero amount should return error",
			fields: fields{
				trfRepo: mock.TransferRepository{
					OnWithinTransaction: func(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error) {
						return txFunc(ctx)
					},
				},
				accRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						if id == "uuid-1" {
							return &model.Account{Balance: 0}, nil
						}
						if id == "uuid-2" {
							return &model.Account{Balance: 10}, nil
						}

						return nil, errors.New("account not found")
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				transferInput: TransferCreateInput{
					AccountOriginID:      "uuid-1",
					AccountDestinationID: "uuid-2",
					Amount:               0,
				},
			},
			want:    nil,
			wantErr: ErrTransferAmountNotPositive,
		},
		{
			name: "origin account balance zero should return error",
			fields: fields{
				trfRepo: mock.TransferRepository{
					OnWithinTransaction: func(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error) {
						return txFunc(ctx)
					},
				},
				accRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						if id == "uuid-1" {
							return &model.Account{Balance: 0}, nil
						}
						if id == "uuid-2" {
							return &model.Account{Balance: 10}, nil
						}

						return nil, errors.New("account not found")
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				transferInput: TransferCreateInput{
					AccountOriginID:      "uuid-1",
					AccountDestinationID: "uuid-2",
					Amount:               1,
				},
			},
			want:    nil,
			wantErr: ErrAccountCurrentBalanceInsufficient,
		},
		{
			name: "origin account balance less than amount should return error",
			fields: fields{
				trfRepo: mock.TransferRepository{
					OnWithinTransaction: func(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error) {
						return txFunc(ctx)
					},
				},
				accRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						if id == "uuid-1" {
							return &model.Account{Balance: 10}, nil
						}
						if id == "uuid-2" {
							return &model.Account{Balance: 10}, nil
						}

						return nil, errors.New("account not found")
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				transferInput: TransferCreateInput{
					AccountOriginID:      "uuid-1",
					AccountDestinationID: "uuid-2",
					Amount:               10.01,
				},
			},
			want:    nil,
			wantErr: ErrAccountCurrentBalanceInsufficient,
		},
		{
			name: "origin account balance not found should return error",
			fields: fields{
				trfRepo: mock.TransferRepository{
					OnWithinTransaction: func(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error) {
						return txFunc(ctx)
					},
				},
				accRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						if id == "uuid-2" {
							return &model.Account{Balance: 1000}, nil
						}

						return nil, repository.ErrAccountNotFound
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				transferInput: TransferCreateInput{
					AccountOriginID:      "uuid-1",
					AccountDestinationID: "uuid-2",
					Amount:               1,
				},
			},
			want:    nil,
			wantErr: repository.ErrAccountNotFound,
		},
		{
			name: "destination account balance not found should return error",
			fields: fields{
				trfRepo: mock.TransferRepository{
					OnWithinTransaction: func(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error) {
						return txFunc(ctx)
					},
				},
				accRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						if id == "uuid-1" {
							return &model.Account{Balance: 1000}, nil
						}

						return nil, repository.ErrAccountNotFound
					},
					OnUpdateBalance: func(ctx context.Context, id model.AccountID, balance model.Money) error {
						return nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				transferInput: TransferCreateInput{
					AccountOriginID:      "uuid-1",
					AccountDestinationID: "uuid-2",
					Amount:               1,
				},
			},
			want:    nil,
			wantErr: repository.ErrAccountNotFound,
		},
		{
			name: "repo create transfer error should return error",
			fields: fields{
				trfRepo: mock.TransferRepository{
					OnWithinTransaction: func(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error) {
						return txFunc(ctx)
					},
					OnCreate: func(ctx context.Context, transfer *model.Transfer) error {
						return errors.New("any error")
					},
				},
				accRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						if id == "uuid-1" {
							return &model.Account{Balance: 1000}, nil
						}
						if id == "uuid-2" {
							return &model.Account{Balance: 1000}, nil
						}

						return nil, repository.ErrAccountNotFound
					},
					OnUpdateBalance: func(ctx context.Context, id model.AccountID, balance model.Money) error {
						return nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				transferInput: TransferCreateInput{
					AccountOriginID:      "uuid-1",
					AccountDestinationID: "uuid-2",
					Amount:               1,
				},
			},
			want:    nil,
			wantErr: ErrTransferCreate,
		},
		{
			name: "success",
			fields: fields{
				trfRepo: mock.TransferRepository{
					OnWithinTransaction: func(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error) {
						return txFunc(ctx)
					},
					OnCreate: func(ctx context.Context, transfer *model.Transfer) error {
						return nil
					},
				},
				accRepo: mock.AccountRepository{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*model.Account, error) {
						if id == "uuid-1" {
							return &model.Account{Balance: 1000}, nil
						}
						if id == "uuid-2" {
							return &model.Account{Balance: 1000}, nil
						}

						return nil, repository.ErrAccountNotFound
					},
					OnUpdateBalance: func(ctx context.Context, id model.AccountID, balance model.Money) error {
						return nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				transferInput: TransferCreateInput{
					AccountOriginID:      "uuid-1",
					AccountDestinationID: "uuid-2",
					Amount:               1.99,
				},
			},
			want: &TransferCreateOutput{
				AccountOriginID:      "uuid-1",
				AccountDestinationID: "uuid-2",
				Amount:               1.99,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trfUC := NewTransferUseCase(tt.fields.trfRepo, tt.fields.accRepo)

			got, err := trfUC.Create(tt.args.ctx, tt.args.transferInput)
			if err != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				if len(got.ID) < 1 {
					t.Errorf("Create() got = %v, want ID generated", got)
				}

				if got.CreatedAt.Before(time.Now().Add(-5 * time.Second)) {
					t.Errorf("Create() got = %v, want CreatedAt in the last 5 seconds", got)
				}

				// clean the generated fields to reflect compare others
				got.ID = ""
				got.CreatedAt = time.Time{}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Create() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
