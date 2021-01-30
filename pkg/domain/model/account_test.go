package model

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
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
