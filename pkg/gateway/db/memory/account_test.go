package memory

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"reflect"
	"testing"
	"time"
)

func TestAccountRepository_Create(t *testing.T) {
	t.Parallel()

	backgroundCtx := context.Background()

	type fields struct {
		accounts []model.Account
	}
	type args struct {
		ctx     context.Context
		account *model.Account
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantRowsCount int
		wantErr       bool
	}{
		{
			name: "empty db success",
			fields: fields{
				accounts: []model.Account{},
			},
			args: args{
				ctx: backgroundCtx,
				account: &model.Account{
					ID:        "uuid-1",
					Name:      "Name 1",
					CPF:       "12345678911",
					Secret:    "whatever",
					Balance:   10,
					CreatedAt: time.Time{},
				},
			},
			wantRowsCount: 1,
			wantErr:       false,
		},
		{
			name: "success with existing account",
			fields: fields{
				accounts: []model.Account{
					{
						ID:        "uuid-1",
						Name:      "Name 1",
						CPF:       "12345678911",
						Secret:    "whatever",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
			},
			args: args{
				account: &model.Account{
					ID:        "uuid-2",
					Name:      "Name 3",
					CPF:       "12345678912",
					Secret:    "whatever2",
					Balance:   102,
					CreatedAt: time.Time{},
				},
			},
			wantRowsCount: 2,
			wantErr:       false,
		},
		{
			name: "existing account id error",
			fields: fields{
				accounts: []model.Account{
					{
						ID:        "uuid-1",
						Name:      "Name 1",
						CPF:       "12345678911",
						Secret:    "whatever",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
			},
			args: args{
				account: &model.Account{
					ID:        "uuid-1",
					Name:      "Name 3",
					CPF:       "12345678912",
					Secret:    "whatever2",
					Balance:   102,
					CreatedAt: time.Time{},
				},
			},
			wantRowsCount: 1,
			wantErr:       true,
		},
		{
			name: "existing account CPF error",
			fields: fields{
				accounts: []model.Account{
					{
						ID:        "uuid-1",
						Name:      "Name 1",
						CPF:       "12345678911",
						Secret:    "whatever",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
			},
			args: args{
				account: &model.Account{
					ID:        "uuid-9",
					Name:      "Name 3",
					CPF:       "12345678911",
					Secret:    "whatever2",
					Balance:   102,
					CreatedAt: time.Time{},
				},
			},
			wantRowsCount: 1,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewAccountRepository(tt.fields.accounts...)
			if err := repo.Create(tt.args.ctx, tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(repo.accountsByIDMap) != tt.wantRowsCount {
				t.Errorf("Create() accountsByIDMap.count = %v, wantRowsCount %v", len(repo.accountsByIDMap), tt.wantRowsCount)
			}

			if !tt.wantErr && !reflect.DeepEqual(repo.accountsByIDMap[tt.args.account.ID], *tt.args.account) {
				t.Errorf("Create() accountsByIDMap.saved = %v, want %v", repo.accountsByIDMap[tt.args.account.ID], tt.args.account)
			}
		})
	}
}

func TestAccountRepository_ExistsByCPF(t *testing.T) {
	t.Parallel()

	backgroundCtx := context.Background()

	type fields struct {
		accounts []model.Account
	}
	type args struct {
		ctx context.Context
		cpf model.CPF
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty db success",
			fields: fields{
				accounts: []model.Account{},
			},
			args: args{
				ctx: backgroundCtx,
				cpf: "12345678911",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success with existing account",
			fields: fields{
				accounts: []model.Account{
					{
						ID:        "uuid-1",
						Name:      "Name 1",
						CPF:       "12345678911",
						Secret:    "whatever",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				cpf: "12345678912",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "existing account CPF error",
			fields: fields{
				accounts: []model.Account{
					{
						ID:        "uuid-1",
						Name:      "Name 1",
						CPF:       "12345678911",
						Secret:    "whatever",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
			},
			args: args{
				ctx: backgroundCtx,
				cpf: "12345678911",
			},
			want:    true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewAccountRepository(tt.fields.accounts...)
			got, err := repo.ExistsByCPF(tt.args.ctx, tt.args.cpf)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExistsByCPF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExistsByCPF() got = %v, want %v", got, tt.want)
			}
		})
	}
}
