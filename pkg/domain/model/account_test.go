package model

import (
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
	"time"
)

func TestAccount_HashSecret(t *testing.T) {
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
			name: "",
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
		// TODO add more tests (formatted cpf, negative balance, positive balance, empty fields)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAccount(tt.args.name, tt.args.cpf, tt.args.secret, tt.args.balance)
			// TODO check if id was generated
			// TODO check if createdAt was generated
			got.ID = ""
			got.CreatedAt = time.Time{}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}