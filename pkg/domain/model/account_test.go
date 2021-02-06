package model

import (
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
	"time"
)

func TestAccount_HashSecret(t *testing.T) {
	t.Parallel()

	type fields struct {
		Secret string
	}
	tests := []struct {
		name           string
		fields         fields
		wantHashSecret string
		wantErr        bool
	}{
		{
			name: "123",
			fields: fields{
				Secret: "123",
			},
			wantErr: false,
		},
		{
			name: "empty",
			fields: fields{
				Secret: "",
			},
			wantErr: false,
		},
		{
			name: "abcABC !@#$%123456",
			fields: fields{
				Secret: "abcABC !@#$%123456",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Secret: tt.fields.Secret,
			}
			if err := a.HashSecret(); (err != nil) != tt.wantErr {
				t.Errorf("HashSecret() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := bcrypt.CompareHashAndPassword([]byte(a.Secret), []byte(tt.fields.Secret)); err != nil {
				t.Errorf("HashSecret() compareHashAndPassword, secret %v, error %v", a.Secret, err)
			}
		})
	}
}

func TestNewAccount(t *testing.T) {
	t.Parallel()

	type args struct {
		name    string
		cpf     string
		secret  string
		balance float64
	}
	tests := []struct {
		name string
		args args
		want *Account
	}{
		{
			name: "success",
			args: args{
				name:    "Bart Simpson",
				cpf:     "12345678911",
				secret:  "123456",
				balance: 0,
			},
			want: &Account{
				ID:        "",
				Name:      "Bart Simpson",
				CPF:       "12345678911",
				Secret:    "123456",
				Balance:   0,
				CreatedAt: time.Time{},
			},
		},
		{
			name: "formatted CPF should return non-formatted",
			args: args{
				name:    "Bart Simpson",
				cpf:     "123.456.789-11",
				secret:  "123456",
				balance: 0,
			},
			want: &Account{
				ID:        "",
				Name:      "Bart Simpson",
				CPF:       "12345678911",
				Secret:    "123456",
				Balance:   0,
				CreatedAt: time.Time{},
			},
		},
		{
			name: "name with lead/trailing spaces should trim",
			args: args{
				name:    "  Bart Simpson  ",
				cpf:     "12345678911",
				secret:  "123456",
				balance: 0,
			},
			want: &Account{
				ID:        "",
				Name:      "Bart Simpson",
				CPF:       "12345678911",
				Secret:    "123456",
				Balance:   0,
				CreatedAt: time.Time{},
			},
		},
		{
			name: "negative balance should be OK",
			args: args{
				name:    "Bart Simpson",
				cpf:     "12345678911",
				secret:  "123456",
				balance: -1.9,
			},
			want: &Account{
				ID:        "",
				Name:      "Bart Simpson",
				CPF:       "12345678911",
				Secret:    "123456",
				Balance:   -190,
				CreatedAt: time.Time{},
			},
		},
		{
			name: "positive balance should be OK",
			args: args{
				name:    "Bart Simpson",
				cpf:     "12345678911",
				secret:  "123456",
				balance: 1.9,
			},
			want: &Account{
				ID:        "",
				Name:      "Bart Simpson",
				CPF:       "12345678911",
				Secret:    "123456",
				Balance:   190,
				CreatedAt: time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAccount(tt.args.name, tt.args.cpf, tt.args.secret, tt.args.balance)

			if len(got.ID) <= 0 {
				t.Errorf("NewAccount() = %v, ID should not be empty", got)
			}
			got.ID = ""

			if got.CreatedAt.Before(time.Now().Add(-5 * time.Second)) {
				t.Errorf("NewAccount() got = %v, want CreatedAt in the last 5 seconds", got)
			}
			got.CreatedAt = time.Time{}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_CompareSecrets(t *testing.T) {
	t.Parallel()

	type fields struct {
		Secret string
	}
	type args struct {
		secret string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Secret: func() string {
					hashedSecret, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
					return string(hashedSecret)
				}(),
			},
			args: args{
				"secret",
			},
			wantErr: false,
		},
		{
			name: "wrong secret should error",
			fields: fields{
				Secret: func() string {
					hashedSecret, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
					return string(hashedSecret)
				}(),
			},
			args: args{
				"wrong",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Secret: tt.fields.Secret,
			}
			if err := a.CompareSecrets(tt.args.secret); (err != nil) != tt.wantErr {
				t.Errorf("CompareSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
