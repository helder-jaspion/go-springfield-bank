package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"reflect"
	"testing"
	"time"
)

func Test_transferRepository_Create(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx      context.Context
		transfer *model.Transfer
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
			name: "should return err when origin account not exists",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				transfer: &model.Transfer{
					ID:                   model.NewTransferID(),
					AccountOriginID:      model.NewAccountID(),
					AccountDestinationID: model.NewAccountID(),
					Amount:               10,
					CreatedAt:            time.Now().Round(time.Microsecond),
				},
			},
			wantErr: true,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
					string(args.transfer.AccountDestinationID),
					"destination",
					"00000000002",
					"any secret")
				if err != nil {
					t.Errorf("Create() error on runBefore = %v", err)
				}
			},
			check: func(args args) {},
		},
		{
			name: "should return err when destination account not exists",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				transfer: &model.Transfer{
					ID:                   model.NewTransferID(),
					AccountOriginID:      model.NewAccountID(),
					AccountDestinationID: model.NewAccountID(),
					Amount:               10,
					CreatedAt:            time.Now().Round(time.Microsecond),
				},
			},
			wantErr: true,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
					string(args.transfer.AccountOriginID),
					"origin",
					"00000000001",
					"any secret")
				if err != nil {
					t.Errorf("Create() error on runBefore = %v", err)
				}
			},
			check: func(args args) {},
		},
		{
			name: "should success",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				transfer: &model.Transfer{
					ID:                   model.NewTransferID(),
					AccountOriginID:      model.NewAccountID(),
					AccountDestinationID: model.NewAccountID(),
					Amount:               10,
					CreatedAt:            time.Now().Round(time.Microsecond),
				},
			},
			wantErr: false,
			runBefore: func(args args) {
				truncateDatabase(t)

				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
					string(args.transfer.AccountOriginID),
					"origin",
					"00000000001",
					"any secret")
				if err != nil {
					t.Errorf("Create() error on runBefore = %v", err)
				}

				_, err = testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
					string(args.transfer.AccountDestinationID),
					"destination",
					"00000000002",
					"any secret")
				if err != nil {
					t.Errorf("Create() error on runBefore = %v", err)
				}
			},
			check: func(args args) {
				var got model.Transfer
				err := testDbPool.QueryRow(backgroundCtx, "SELECT id, account_origin_id, account_destination_id, amount, created_at FROM transfers WHERE id = $1", string(args.transfer.ID)).
					Scan(&got.ID, &got.AccountOriginID, &got.AccountDestinationID, &got.Amount, &got.CreatedAt)
				if err != nil {
					t.Errorf("Create() error = %v, wantErr %v", err, false)
				}
				if !reflect.DeepEqual(got, *args.transfer) {
					t.Errorf("Create() got = %v, want %v", got, *args.transfer)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			trfRepo := NewTransferRepository(tt.fields.db)
			if err := trfRepo.Create(tt.args.ctx, tt.args.transfer); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			tt.check(tt.args)
		})
	}
}

func Test_transferRepository_Fetch(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx       context.Context
		accountID model.AccountID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      func(args) []model.Transfer
		wantErr   bool
		runBefore func(args, []model.Transfer)
	}{
		{
			name: "should return empty if db empty",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx:       backgroundCtx,
				accountID: model.AccountID(uuid.NewString()),
			},
			want: func(args args) []model.Transfer {
				return []model.Transfer{}
			},
			wantErr: false,
			runBefore: func(args args, _ []model.Transfer) {
				truncateDatabase(t)
			},
		},
		{
			name: "should return one when own the origin account",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx:       backgroundCtx,
				accountID: model.NewAccountID(),
			},
			want: func(args args) []model.Transfer {
				return []model.Transfer{
					{
						ID:                   model.NewTransferID(),
						AccountOriginID:      args.accountID,
						AccountDestinationID: model.NewAccountID(),
						Amount:               123,
						CreatedAt:            time.Now().Round(time.Microsecond),
					},
				}
			},
			wantErr: false,
			runBefore: func(args args, values []model.Transfer) {
				truncateDatabase(t)

				for i, v := range values {
					_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
						string(v.AccountOriginID),
						"origin",
						fmt.Sprintf("%011d", i+1),
						"any secret")
					if err != nil {
						t.Errorf("Create() error on runBefore = %v", err)
					}

					_, err = testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
						string(v.AccountDestinationID),
						"destination",
						fmt.Sprintf("%011d", (i+1)*2),
						"any secret")
					if err != nil {
						t.Errorf("Create() error on runBefore = %v", err)
					}

					_, err = testDbPool.Exec(backgroundCtx, "INSERT INTO transfers (id, account_origin_id, account_destination_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
						string(v.ID),
						string(v.AccountOriginID),
						string(v.AccountDestinationID),
						v.Amount,
						v.CreatedAt)
					if err != nil {
						t.Errorf("Fetch() error on runBefore = %v", err)
					}
				}
			},
		},
		{
			name: "should return one when own the destination account",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx:       backgroundCtx,
				accountID: model.NewAccountID(),
			},
			want: func(args args) []model.Transfer {
				return []model.Transfer{
					{
						ID:                   model.NewTransferID(),
						AccountOriginID:      model.NewAccountID(),
						AccountDestinationID: args.accountID,
						Amount:               123,
						CreatedAt:            time.Now().Round(time.Microsecond),
					},
				}
			},
			wantErr: false,
			runBefore: func(args args, values []model.Transfer) {
				truncateDatabase(t)

				for i, v := range values {
					_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
						string(v.AccountOriginID),
						"origin",
						fmt.Sprintf("%011d", i+1),
						"any secret")
					if err != nil {
						t.Errorf("Create() error on runBefore = %v", err)
					}

					_, err = testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
						string(v.AccountDestinationID),
						"destination",
						fmt.Sprintf("%011d", (i+1)*2),
						"any secret")
					if err != nil {
						t.Errorf("Create() error on runBefore = %v", err)
					}

					_, err = testDbPool.Exec(backgroundCtx, "INSERT INTO transfers (id, account_origin_id, account_destination_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
						string(v.ID),
						string(v.AccountOriginID),
						string(v.AccountDestinationID),
						v.Amount,
						v.CreatedAt)
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
				ctx:       backgroundCtx,
				accountID: model.NewAccountID(),
			},
			want: func(args args) []model.Transfer {
				return []model.Transfer{
					{
						ID:                   model.NewTransferID(),
						AccountOriginID:      model.NewAccountID(),
						AccountDestinationID: args.accountID,
						Amount:               123,
						CreatedAt:            time.Now().Round(time.Microsecond),
					},
					{
						ID:                   model.NewTransferID(),
						AccountOriginID:      args.accountID,
						AccountDestinationID: model.NewAccountID(),
						Amount:               111,
						CreatedAt:            time.Now().Add(-1 * time.Minute).Round(time.Microsecond),
					},
				}
			},
			wantErr: false,
			runBefore: func(args args, values []model.Transfer) {
				truncateDatabase(t)

				for i, v := range values {
					_, _ = testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
						string(v.AccountOriginID),
						"origin",
						fmt.Sprintf("%011d", i+1),
						"any secret")

					_, _ = testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
						string(v.AccountDestinationID),
						"destination",
						fmt.Sprintf("%011d", (i+1)*2),
						"any secret")

					_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO transfers (id, account_origin_id, account_destination_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
						string(v.ID),
						string(v.AccountOriginID),
						string(v.AccountDestinationID),
						v.Amount,
						v.CreatedAt)
					if err != nil {
						t.Errorf("Fetch() error on runBefore = %v", err)
					}
				}
			},
		},
		{
			name: "should return empty when not own transfers",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx:       backgroundCtx,
				accountID: model.NewAccountID(),
			},
			want: func(args args) []model.Transfer {
				return []model.Transfer{}
			},
			wantErr: false,
			runBefore: func(args args, values []model.Transfer) {
				truncateDatabase(t)

				originUUID, destinationUUID := uuid.NewString(), uuid.NewString()
				_, err := testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
					originUUID,
					"origin",
					"00011100011",
					"any secret")
				if err != nil {
					t.Errorf("Create() error on runBefore = %v", err)
				}

				_, err = testDbPool.Exec(backgroundCtx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
					destinationUUID,
					"destination",
					"00011100022",
					"any secret")
				if err != nil {
					t.Errorf("Create() error on runBefore = %v", err)
				}

				_, err = testDbPool.Exec(backgroundCtx, "INSERT INTO transfers (id, account_origin_id, account_destination_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
					uuid.NewString(),
					originUUID,
					destinationUUID,
					1000,
					time.Now())
				if err != nil {
					t.Errorf("Fetch() error on runBefore = %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want(tt.args)
			if tt.runBefore != nil {
				tt.runBefore(tt.args, want)
			}

			trfRepo := NewTransferRepository(tt.fields.db)
			got, err := trfRepo.Fetch(tt.args.ctx, tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Fetch() got = %v, want %v", got, want)
			}
		})
	}
}

func Test_transferRepository_WithinTransaction(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx    context.Context
		txFunc func(context.Context) (interface{}, error)
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantData  interface{}
		wantErr   bool
		runBefore func()
		check     func()
	}{
		{
			name: "should rollback",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				txFunc: func(ctx context.Context) (interface{}, error) {
					originUUID, destinationUUID := uuid.NewString(), uuid.NewString()
					_, err := getConnFromCtx(ctx, testDbPool).Exec(ctx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
						originUUID,
						"origin",
						"00011100011",
						"any secret")
					if err != nil {
						t.Errorf("WithinTransaction() error on runBefore = %v", err)
					}

					_, err = getConnFromCtx(ctx, testDbPool).Exec(backgroundCtx, "INSERT INTO transfers (id, account_origin_id, account_destination_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
						uuid.NewString(),
						originUUID,
						destinationUUID,
						1000,
						time.Now())
					if err == nil {
						t.Error("WithinTransaction() should error on runBefore")
					}

					return nil, err
				},
			},
			wantData: nil,
			wantErr:  true,
			runBefore: func() {
				truncateDatabase(t)
			},
			check: func() {
				accountCount := 0
				err := testDbPool.QueryRow(backgroundCtx, "SELECT COUNT(id) FROM accounts").Scan(&accountCount)
				if err != nil {
					t.Errorf("WithinTransaction() error = %v, wantErr %v", err, false)
				}
				if accountCount != 0 {
					t.Error("WithinTransaction() accountCount should be 0")
				}

				transferCount := 0
				err = testDbPool.QueryRow(backgroundCtx, "SELECT COUNT(id) FROM transfers").Scan(&transferCount)
				if err != nil {
					t.Errorf("WithinTransaction() error = %v, wantErr %v", err, false)
				}
				if transferCount != 0 {
					t.Error("WithinTransaction() transferCount should be 0")
				}
			},
		},
		{
			name: "should commit",
			fields: fields{
				db: testDbPool,
			},
			args: args{
				ctx: backgroundCtx,
				txFunc: func(ctx context.Context) (interface{}, error) {
					originUUID, destinationUUID := uuid.NewString(), uuid.NewString()
					_, err := getConnFromCtx(ctx, testDbPool).Exec(ctx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
						originUUID,
						"origin",
						"00011100011",
						"any secret")
					if err != nil {
						t.Errorf("WithinTransaction() error on runBefore = %v", err)
					}

					_, err = getConnFromCtx(ctx, testDbPool).Exec(ctx, "INSERT INTO accounts (id, name, cpf, secret) VALUES ($1, $2, $3, $4)",
						destinationUUID,
						"destination",
						"00011100022",
						"any secret")
					if err != nil {
						t.Errorf("WithinTransaction() error on runBefore = %v", err)
					}

					_, err = getConnFromCtx(ctx, testDbPool).Exec(backgroundCtx, "INSERT INTO transfers (id, account_origin_id, account_destination_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
						uuid.NewString(),
						originUUID,
						destinationUUID,
						1000,
						time.Now())
					if err != nil {
						t.Errorf("WithinTransaction() error on runBefore = %v", err)
					}

					return nil, err
				},
			},
			wantData: nil,
			wantErr:  false,
			runBefore: func() {
				truncateDatabase(t)
			},
			check: func() {
				accountCount := 0
				err := testDbPool.QueryRow(backgroundCtx, "SELECT COUNT(id) FROM accounts").Scan(&accountCount)
				if err != nil {
					t.Errorf("WithinTransaction() error = %v, wantErr %v", err, false)
				}
				if accountCount != 2 {
					t.Error("WithinTransaction() accountCount should be 2")
				}

				transferCount := 0
				err = testDbPool.QueryRow(backgroundCtx, "SELECT COUNT(id) FROM transfers").Scan(&transferCount)
				if err != nil {
					t.Errorf("WithinTransaction() error = %v, wantErr %v", err, false)
				}
				if transferCount != 1 {
					t.Error("WithinTransaction() transferCount should be 1")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore()
			}

			trfRepo := NewTransferRepository(tt.fields.db)
			gotData, err := trfRepo.WithinTransaction(tt.args.ctx, tt.args.txFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithinTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("WithinTransaction() gotData = %v, want %v", gotData, tt.wantData)
			}

			tt.check()
		})
	}
}
