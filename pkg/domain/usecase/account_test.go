package usecase

import (
	"context"
	"errors"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestAccountCreateInput_Validate(t *testing.T) {
	type fields struct {
		Name    string
		CPF     string
		Secret  string
		Balance float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "Empty name should return error",
			fields: fields{
				Name:    "",
				CPF:     "",
				Secret:  "",
				Balance: 0,
			},
			wantErr: ErrAccountNameWrongLength,
		},
		{
			name: "one char name should return error",
			fields: fields{
				Name:    "A",
				CPF:     "",
				Secret:  "",
				Balance: 0,
			},
			wantErr: ErrAccountNameWrongLength,
		},
		{
			name: "101 chars name should return error",
			fields: fields{
				Name:    strings.Repeat("A", 101),
				CPF:     "",
				Secret:  "",
				Balance: 0,
			},
			wantErr: ErrAccountNameWrongLength,
		},
		{
			name: "empty CPF should return error",
			fields: fields{
				Name:    "Jon Snow",
				CPF:     "",
				Secret:  "",
				Balance: 0,
			},
			wantErr: ErrAccountCPFInvalid,
		},
		{
			name: "10 digits CPF should return error",
			fields: fields{
				Name:    "Jon Snow",
				CPF:     "1234567890",
				Secret:  "",
				Balance: 0,
			},
			wantErr: ErrAccountCPFInvalid,
		},
		{
			name: "12 digits CPF should return error",
			fields: fields{
				Name:    "Jon Snow",
				CPF:     "123456789012",
				Secret:  "",
				Balance: 0,
			},
			wantErr: ErrAccountCPFInvalid,
		},
		{
			name: "invalid CPF should return error",
			fields: fields{
				Name:    "Jon Snow",
				CPF:     "12345678901",
				Secret:  "",
				Balance: 0,
			},
			wantErr: ErrAccountCPFInvalid,
		},
		{
			name: "empty secret should return error",
			fields: fields{
				Name:    "Jon Snow",
				CPF:     "599.513.320-99",
				Secret:  "",
				Balance: 0,
			},
			wantErr: ErrAccountSecretWrongLength,
		},
		{
			name: "5 chars secret should return error",
			fields: fields{
				Name:    "Jon Snow",
				CPF:     "599.513.320-99",
				Secret:  "12345",
				Balance: 0,
			},
			wantErr: ErrAccountSecretWrongLength,
		},
		{
			name: "101 chars secret should return error",
			fields: fields{
				Name:    "Jon Snow",
				CPF:     "599.513.320-99",
				Secret:  strings.Repeat("A", 101),
				Balance: 0,
			},
			wantErr: ErrAccountSecretWrongLength,
		},
		{
			name: "negative balance should return error",
			fields: fields{
				Name:    "Jon Snow",
				CPF:     "599.513.320-99",
				Secret:  "IAmNotSnow",
				Balance: -1,
			},
			wantErr: ErrAccountBalanceNegative,
		},
		{
			name: "positive decimal balance should success",
			fields: fields{
				Name:    "Jon Snow",
				CPF:     "599.513.320-99",
				Secret:  "IAmNotSnow",
				Balance: 0.1,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := AccountCreateInput{
				Name:    tt.fields.Name,
				CPF:     tt.fields.CPF,
				Secret:  tt.fields.Secret,
				Balance: tt.fields.Balance,
			}
			err := input.Validate()
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_accountUseCase_Create(t *testing.T) {
	t.Parallel()

	backgroundCtx := context.Background()

	type fields struct {
		accountRepo repository.AccountRepository
	}
	type args struct {
		ctx          context.Context
		accountInput AccountCreateInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AccountCreateOutput
		wantErr bool
	}{
		{
			name: "repo create error should return error",
			fields: fields{
				accountRepo: repository.AccountRepositoryMock{
					OnCreate: func(ctx context.Context, account *model.Account) error {
						return errors.New("any database error")
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				accountInput: AccountCreateInput{
					Name:    "Jon Snow",
					CPF:     "599.513.320-99",
					Secret:  "IAmNotSnow",
					Balance: 0,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "input name empty should return error",
			fields: fields{
				accountRepo: repository.AccountRepositoryMock{},
			},
			args: args{
				ctx: backgroundCtx,
				accountInput: AccountCreateInput{
					Name:    "",
					CPF:     "599.513.320-99",
					Secret:  "IAmNotSnow",
					Balance: 0,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unformatted CPF should return formatted",
			fields: fields{
				accountRepo: repository.AccountRepositoryMock{
					OnCreate: func(ctx context.Context, account *model.Account) error {
						return nil
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				accountInput: AccountCreateInput{
					Name:    "Jon Snow",
					CPF:     "59951332099",
					Secret:  "IAmNotSnow",
					Balance: 0,
				},
			},
			want: &AccountCreateOutput{
				Name:    "Jon Snow",
				CPF:     "599.513.320-99",
				Balance: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountUC := NewAccountUseCase(tt.fields.accountRepo)

			got, err := accountUC.Create(tt.args.ctx, tt.args.accountInput)
			if (err != nil) != tt.wantErr {
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
