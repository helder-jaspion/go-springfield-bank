package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

func Test_authUseCase_Authorize(t *testing.T) {
	t.Parallel()

	backgroundCtx := context.Background()

	type fields struct {
		secretKey      string
		accessTokenDur time.Duration
		accRepo        repository.AccountRepository
	}
	type args struct {
		ctx         context.Context
		accessToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *jwt.RegisteredClaims
		wantErr error
	}{
		{
			name: "expired token should return invalid access token",
			fields: fields{
				secretKey:      "whatever",
				accessTokenDur: 1 * time.Minute,
			},
			args: args{
				ctx:         backgroundCtx,
				accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTIzMDMyMjQsImlhdCI6MTYxMjMwMjMyNCwic3ViIjoiNDQxOTc0NjQtMDJiZS00OWMyLWIwZmQtNjYwZmFlZDczNWJkIn0.-Iz5yzIMwLeYz0kFknwHwwBAngaBdTyyMK826WkUjE4",
			},
			want:    nil,
			wantErr: ErrAuthInvalidAccessToken,
		},
		{
			name: "wrong signing method token should return invalid access token",
			fields: fields{
				secretKey:      "whatever",
				accessTokenDur: 1 * time.Minute,
			},
			args: args{
				ctx:         backgroundCtx,
				accessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.POstGetfAytaZS82wHcjoTyoqhMyxXiWdR7Nn7A29DNSl0EiXLdwJ6xC6AfgZWF1bOsS_TuYI3OG85AmiExREkrS6tDfTQ2B3WXlrr-wp5AokiRbz3_oB4OxG-W9KcEEbDRcZc0nH3L7LzYptiy1PtAylQGxHTWZXtGz4ht0bAecBgmpdgXMguEIcoqPJ1n3pIWk_dUZegpqx0Lka21H6XxUTxiy8OcaarA8zdnPUnV6AmNP3ecFawIFYdvJB_cm-GvpCSbr8G8y_Mllj8f4x9nBH8pQux89_6gUY618iYv7tuPWBFfEbLxtF2pZS6YC1aSfLQxeNe8djT9YjpvRZA",
			},
			want:    nil,
			wantErr: ErrAuthInvalidAccessToken,
		},
		{
			name: "valid access token",
			fields: fields{
				secretKey:      "whatever",
				accessTokenDur: 1 * time.Minute,
			},
			args: args{
				ctx:         backgroundCtx,
				accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2MTIzMDQzMTcsInN1YiI6IjQ0MTk3NDY0LTAyYmUtNDljMi1iMGZkLTY2MGZhZWQ3MzViZCJ9.tpCm8rPCsWG9exmw5_Ic9dGohNo2Q5S_PrCZirH8x-w",
			},
			want: &jwt.RegisteredClaims{
				IssuedAt: jwt.NewNumericDate(time.Date(2021, time.February, 2, 22, 18, 37, 0, time.UTC)),
				Subject:  "44197464-02be-49c2-b0fd-660faed735bd",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authUC := NewAuthUseCase(
				tt.fields.secretKey,
				tt.fields.accessTokenDur,
				tt.fields.accRepo,
			)

			got, err := authUC.Authorize(tt.args.ctx, tt.args.accessToken)
			if err != tt.wantErr {
				t.Errorf("Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				if got.IssuedAt.UTC() != tt.want.IssuedAt.UTC() {
					t.Errorf("Authorize() got = %v, want %v", got, tt.want)
					return
				}

				got.IssuedAt = tt.want.IssuedAt
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Authorize() got = %v, want %v", got, tt.want)
			}
		})
	}
}
