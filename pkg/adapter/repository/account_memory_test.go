package repository

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestAccountMemoryRepository_Create(t *testing.T) {
	t.Parallel()

	backgroundCtx := context.Background()

	type fields struct {
		accounts      map[string]model.Account
		accountsByCPF map[string]model.Account
		lock          *sync.RWMutex
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
				accounts:      map[string]model.Account{},
				accountsByCPF: map[string]model.Account{},
				lock:          &sync.RWMutex{},
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
				accounts: map[string]model.Account{
					"uuid-1": {
						ID:        "uuid-1",
						Name:      "Name 1",
						CPF:       "12345678911",
						Secret:    "whatever",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
				accountsByCPF: map[string]model.Account{
					"12345678911": {
						ID:        "uuid-1",
						Name:      "Name 1",
						CPF:       "12345678911",
						Secret:    "whatever",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
				lock: &sync.RWMutex{},
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
				accounts: map[string]model.Account{
					"uuid-1": {
						ID:        "uuid-1",
						Name:      "Name 1",
						CPF:       "12345678911",
						Secret:    "whatever",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
				accountsByCPF: map[string]model.Account{
					"12345678911": {
						ID:        "uuid-1",
						Name:      "Name 1",
						CPF:       "12345678911",
						Secret:    "whatever",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
				lock: &sync.RWMutex{},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := AccountMemoryRepository{
				accounts: tt.fields.accounts,
				lock:     tt.fields.lock,
			}
			if err := repo.Create(tt.args.ctx, tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(repo.accounts) != tt.wantRowsCount {
				t.Errorf("Create() accounts.count = %v, wantRowsCount %v", len(repo.accounts), tt.wantRowsCount)
			}

			if !tt.wantErr && !reflect.DeepEqual(repo.accounts[tt.args.account.ID], *tt.args.account) {
				t.Errorf("Create() accounts.saved = %v, want %v", repo.accounts[tt.args.account.ID], tt.args.account)
			}
		})
	}
}
