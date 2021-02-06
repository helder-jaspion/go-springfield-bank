package redis

import (
	"context"
	"github.com/go-redis/redis"
	"reflect"
	"testing"
	"time"
)

func Test_idempotencyRepository_Get(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []byte
		wantErr   bool
		runBefore func(args)
	}{
		{
			name: "empty cache should return error",
			fields: fields{
				client: testRedisClient,
			},
			args: args{
				ctx: backgroundCtx,
				key: "any-key-1",
			},
			want:    []byte{},
			wantErr: true,
		},
		{
			name: "should success",
			fields: fields{
				client: testRedisClient,
			},
			args: args{
				ctx: backgroundCtx,
				key: "any-key-2",
			},
			runBefore: func(args args) {
				testRedisClient.Set("_IDEMPOTENCY_"+args.key, "value", 10*time.Second)
			},
			want:    []byte("value"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore != nil {
				tt.runBefore(tt.args)
			}

			idpRepo := NewIdempotencyRepository(tt.fields.client)

			got, err := idpRepo.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_idempotencyRepository_Set(t *testing.T) {
	backgroundCtx := context.Background()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		ctx      context.Context
		key      string
		value    []byte
		duration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		check   func(args)
	}{
		{
			name: "should success",
			fields: fields{
				client: testRedisClient,
			},
			args: args{
				ctx:      backgroundCtx,
				key:      "any-key-3",
				value:    []byte("value-3"),
				duration: 10 * time.Second,
			},
			wantErr: false,
			check: func(args args) {
				value, err := testRedisClient.Get("_IDEMPOTENCY_" + args.key).Bytes()
				if err != nil {
					t.Errorf("Set() error = %v, wantErr %v", err, false)
				}

				if !reflect.DeepEqual(value, args.value) {
					t.Errorf("Set() value = %v, wantValue %v", value, args.value)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idpRepo := NewIdempotencyRepository(tt.fields.client)

			if err := idpRepo.Set(tt.args.ctx, tt.args.key, tt.args.value, tt.args.duration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			tt.check(tt.args)
		})
	}
}
