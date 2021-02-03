package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// AccountID represents an Account ID as uuid.
type AccountID string

// NewAccountID returns a new AccountID with value generated by uuid.New().
func NewAccountID() AccountID {
	return AccountID(uuid.New().String())
}

// Account represents a bank account.
type Account struct {
	ID        AccountID
	Name      string
	CPF       CPF
	Secret    string
	Balance   Money
	CreatedAt time.Time
}

// NewAccount returns a new Account filled with the corresponding arguments with generated values for id and createdAt.
func NewAccount(name string, cpf string, secret string, balance float64) *Account {
	return &Account{
		ID:        NewAccountID(),
		Name:      strings.TrimSpace(name),
		CPF:       NewCPF(cpf),
		Secret:    secret,
		Balance:   Float64ToMoney(balance),
		CreatedAt: time.Now(),
	}
}

// HashSecret hashes secret with bcrypt.
func (a *Account) HashSecret() error {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(a.Secret), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Secret = string(hashedSecret)
	return nil
}

// CompareSecrets compare account secret and payload.
func (a *Account) CompareSecrets(secret string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(a.Secret), []byte(secret)); err != nil {
		return err
	}
	return nil
}
