package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"reflect"
	"testing"
	"time"
)

func Test_accountRepository_Create(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx     context.Context
		account *model.Account
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		runBefore func(args)
		check     func(args)
	}{
		{
			name: "should success",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				account: &model.Account{
					ID:        model.AccountID(uuid.NewString()),
					Name:      "Bart Simpson",
					CPF:       "12345678911",
					Secret:    "secret",
					Balance:   100,
					CreatedAt: time.Now().Round(time.Microsecond),
				},
			},
			wantErr: false,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
			check: func(args args) {
				var got model.Account
				err := testDbPool.QueryRow(backgroundCtx, "SELECT id, name, cpf, secret, balance, created_at FROM accounts WHERE id = $1", string(args.account.ID)).
					Scan(&got.ID, &got.Name, &got.CPF, &got.Secret, &got.Balance, &got.CreatedAt)
				if err != nil {
					t.Errorf("Create() error = %v, wantErr %v", err, false)
				}
				if !reflect.DeepEqual(got, *args.account) {
					t.Errorf("Create() got = %v, want %v", got, *args.account)
				}
			},
		},
		{
			name: "existent with the same ID should return error",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				account: &model.Account{
					ID:        model.AccountID(uuid.NewString()),
					Name:      "Bart Simpson",
					CPF:       "12345678911",
					Secret:    "secret",
					Balance:   100,
					CreatedAt: time.Now().Round(time.Microsecond),
				},
			},
			wantErr: true,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
					string(args.account.ID),
					args.account.Name,
					"00000000001",
					args.account.Secret)
				if err != nil {
					t.Errorf("Create() error on runBefore = %v", err)
				}
			},
			check: func(args args) {},
		},
		{
			name: "existent with the same CPF should return error",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				account: &model.Account{
					ID:        model.AccountID(uuid.NewString()),
					Name:      "Bart Simpson",
					CPF:       "12345678911",
					Secret:    "secret",
					Balance:   100,
					CreatedAt: time.Now().Round(time.Microsecond),
				},
			},
			wantErr: true,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
					uuid.NewString(),
					args.account.Name,
					"12345678911",
					args.account.Secret)
				if err != nil {
					t.Errorf("Create() error on runBefore = %v", err)
				}
			},
			check: func(args args) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			accRepo := NewAccountRepository(tt.fields.db)
			if err := accRepo.Create(tt.args.ctx, tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			tt.check(tt.args)
		})
	}
}

func Test_accountRepository_ExistsByCPF(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		cpf model.CPF
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		wantErr   bool
		runBefore func(args)
	}{
		{
			name: "should return false if empty db",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				cpf: "12345678911",
			},
			want:    false,
			wantErr: false,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "should return true if exists",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				cpf: "12345678911",
			},
			want:    true,
			wantErr: false,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
					uuid.NewString(),
					"Bart Simpson",
					"12345678911",
					"secret")
				if err != nil {
					t.Errorf("Create() error on runBefore = %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			accRepo := NewAccountRepository(tt.fields.db)
			got, err := accRepo.ExistsByCPF(tt.args.ctx, tt.args.cpf)
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

func Test_accountRepository_Fetch(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []model.Account
		wantErr   bool
		runBefore func(args, []model.Account)
	}{
		{
			name: "should return empty if db empty",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
			},
			want:    []model.Account{},
			wantErr: false,
			runBefore: func(args args, _ []model.Account) {
				truncateDatabase(t)
			},
		},
		{
			name: "should return one",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
			},
			want: []model.Account{
				{
					ID:        model.AccountID(uuid.NewString()),
					Name:      "Account 001",
					CPF:       "00000000001",
					Secret:    "secret001",
					Balance:   1,
					CreatedAt: time.Now().Round(time.Microsecond),
				},
			},
			wantErr: false,
			runBefore: func(args args, values []model.Account) {
				truncateDatabase(t)

				for _, v := range values {
					_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
						string(v.ID), v.Name, v.CPF, v.Secret, v.Balance, v.CreatedAt)
					if err != nil {
						t.Errorf("Fetch() error on runBefore = %v", err)
					}
				}
			},
		},
		{
			name: "should return two",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
			},
			want: []model.Account{
				{
					ID:        model.AccountID(uuid.NewString()),
					Name:      "Account 001",
					CPF:       "00000000001",
					Secret:    "secret001",
					Balance:   1,
					CreatedAt: time.Now().Round(time.Microsecond),
				},
				{
					ID:        model.AccountID(uuid.NewString()),
					Name:      "Account 002",
					CPF:       "00000000002",
					Secret:    "secret002",
					Balance:   2,
					CreatedAt: time.Now().Round(time.Microsecond),
				},
			},
			wantErr: false,
			runBefore: func(args args, values []model.Account) {
				truncateDatabase(t)

				for _, v := range values {
					_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
						string(v.ID), v.Name, v.CPF, v.Secret, v.Balance, v.CreatedAt)
					if err != nil {
						t.Errorf("Fetch() error on runBefore = %v", err)
					}
				}
			},
		},
		{
			name: "should return two correct sort",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
			},
			want: []model.Account{
				{
					ID:        model.AccountID(uuid.NewString()),
					Name:      "Account 001",
					CPF:       "00000000001",
					Secret:    "secret001",
					Balance:   1,
					CreatedAt: time.Now().Add(-1 * time.Minute).Round(time.Microsecond),
				},
				{
					ID:        model.AccountID(uuid.NewString()),
					Name:      "Account 002",
					CPF:       "00000000002",
					Secret:    "secret002",
					Balance:   2,
					CreatedAt: time.Now().Round(time.Microsecond),
				},
			},
			wantErr: false,
			runBefore: func(args args, values []model.Account) {
				truncateDatabase(t)

				// reverse order
				for i := len(values) - 1; i >= 0; i-- {
					v := values[i]
					_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
						string(v.ID), v.Name, v.CPF, v.Secret, v.Balance, v.CreatedAt)
					if err != nil {
						t.Errorf("Fetch() error on runBefore = %v", err)
					}
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore(tt.args, tt.want)
			}

			accRepo := NewAccountRepository(tt.fields.db)
			got, err := accRepo.Fetch(tt.args.ctx)
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

func Test_accountRepository_GetBalance(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		id  model.AccountID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      func(args) *model.Account
		wantErr   bool
		runBefore func(args)
	}{
		{
			name: "should return error if empty db",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				id:  model.AccountID(uuid.NewString()),
			},
			want: func(args args) *model.Account {
				return nil
			},
			wantErr: true,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "should return error if not found",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				id:  model.AccountID(uuid.NewString()),
			},
			want: func(args args) *model.Account {
				return nil
			},
			wantErr: true,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
					uuid.NewString(), "Any name", "11111111111", "any secret", 0, time.Now())
				if err != nil {
					t.Errorf("Fetch() error on runBefore = %v", err)
				}
			},
		},
		{
			name: "should return success",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				id:  model.AccountID(uuid.NewString()),
			},
			want: func(args args) *model.Account {
				return &model.Account{
					ID:      args.id,
					Balance: 1050,
				}
			},
			wantErr: false,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
					string(args.id), "Any name", "12345678911", "any secret", 1050, time.Now())
				if err != nil {
					t.Errorf("Fetch() error on runBefore = %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			accRepo := NewAccountRepository(tt.fields.db)
			got, err := accRepo.GetBalance(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want := tt.want(tt.args)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("GetBalance() got = %v, want %v", got, want)
			}
		})
	}
}

func Test_accountRepository_GetByCPF(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx   context.Context
		cpf   model.CPF
		genID string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      func(args) *model.Account
		wantErr   bool
		runBefore func(args)
	}{
		{
			name: "should return error if db is empty",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				cpf: "12345678911",
			},
			want: func(args args) *model.Account {
				return nil
			},
			wantErr: true,
			runBefore: func(args args) {
				truncateDatabase(t)
			},
		},
		{
			name: "should return error if not found",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				cpf: "00000000001",
			},
			want: func(args args) *model.Account {
				return nil
			},
			wantErr: true,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
					uuid.NewString(), "Any name", "11111111111", "any secret", 0, time.Now())
				if err != nil {
					t.Errorf("Fetch() error on runBefore = %v", err)
				}
			},
		},
		{
			name: "should return success",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx:   backgroundCtx,
				cpf:   "12345678901",
				genID: uuid.NewString(),
			},
			want: func(args args) *model.Account {
				return &model.Account{
					ID:        model.AccountID(args.genID),
					Name:      "Bart Simpson 001",
					CPF:       "12345678901",
					Secret:    "any secret",
					Balance:   1050,
					CreatedAt: time.Date(2021, 01, 04, 11, 51, 59, 0, time.Local),
				}
			},
			wantErr: false,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
					args.genID, "Bart Simpson 001", "12345678901", "any secret", 1050, time.Date(2021, 01, 04, 11, 51, 59, 0, time.Local))
				if err != nil {
					t.Errorf("Fetch() error on runBefore = %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			accRepo := NewAccountRepository(tt.fields.db)
			got, err := accRepo.GetByCPF(tt.args.ctx, tt.args.cpf)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByCPF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want := tt.want(tt.args)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("GetByCPF() got = %v, want %v", got, want)
			}
		})
	}
}

func Test_accountRepository_UpdateBalance(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx     context.Context
		id      model.AccountID
		balance model.Money
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		runBefore func(args)
		check     func(args)
	}{
		{
			name: "should success",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx:     backgroundCtx,
				id:      model.AccountID(uuid.NewString()),
				balance: 100,
			},
			wantErr: false,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret, balance) VALUES ($1, $2, $3, $4, $5)",
					string(args.id),
					"any name",
					"12345678911",
					"any secret",
					5999)
				if err != nil {
					t.Errorf("UpdateBalance() error on runBefore = %v", err)
				}
			},
			check: func(args args) {
				var got model.Account
				err := testDbPool.QueryRow(backgroundCtx, "SELECT id, name, cpf, secret, balance, created_at FROM accounts WHERE id = $1", string(args.id)).
					Scan(&got.ID, &got.Name, &got.CPF, &got.Secret, &got.Balance, &got.CreatedAt)
				if err != nil {
					t.Errorf("UpdateBalance() error = %v, wantErr %v", err, false)
				}
				if got.Balance != args.balance {
					t.Errorf("UpdateBalance() got balance = %v, want %v", got, args.balance)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			accRepo := NewAccountRepository(tt.fields.db)
			if err := accRepo.UpdateBalance(tt.args.ctx, tt.args.id, tt.args.balance); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBalance() error = %v, wantErr %v", err, tt.wantErr)
			}

			tt.check(tt.args)
		})
	}
}
