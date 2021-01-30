package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// Account represents a bank account.
type Account struct {
	ID        string
	Name      string
	CPF       string
	Secret    string
	Balance   int64
	CreatedAt time.Time
}

// NewAccount returns a new Account filled with the corresponding arguments with generated values for id and createdAt.
func NewAccount(name string, cpf string, secret string, balance float64) *Account {
	return &Account{
		ID:        uuid.New().String(),
		Name:      strings.TrimSpace(name),
		CPF:       cpf,
		Secret:    secret,
		Balance:   int64(balance * 100),
		CreatedAt: time.Now(),
	}
}

// HashSecret hashes the account secret.
func (a *Account) HashSecret() error {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(a.Secret), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Secret = string(hashedSecret)
	return nil
}
