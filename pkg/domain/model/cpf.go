package model

import (
	"github.com/helder-jaspion/go-springfield-bank/pkg/cpfutil"
)

// CPF represents the Brazilian taxpayer ID, know as CPF.
type CPF string

// NewCPF instantiates a new CPF.
// It removes the non-digits from the input, if any
func NewCPF(cpf string) CPF {
	return CPF(cpfutil.Clean(cpf))
}

// String returns a formatted CPF (000.000.000-00)
func (c *CPF) String() string {
	return cpfutil.Format(string(*c))
}
