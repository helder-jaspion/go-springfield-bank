package cpfutil

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Clean returns only-digits version of the cpf.
func Clean(cpf string) string {
	buf := bytes.NewBufferString("")
	for _, r := range cpf {
		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		}
	}

	return buf.String()
}

// Format returns a formatted CPF (000.000.000-00).
func Format(cpf string) string {
	expr, err := regexp.Compile(`^([\d]{3})([\d]{3})([\d]{3})([\d]{2})$`)
	if err != nil {
		return cpf
	}

	return expr.ReplaceAllString(cpf, "$1.$2.$3-$4")
}

// IsValid returns true if it is a valid CPF.
func IsValid(cpf string) bool {
	cpf = Clean(cpf)
	if len(cpf) != 11 {
		return false
	}

	ds := make([]int64, 11)
	allEq := true
	var lastDigit *int64
	for i, v := range strings.Split(cpf, "") {
		c, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return false
		}
		ds[i] = c
		allEq = allEq && (lastDigit == nil || *lastDigit == c)
		lastDigit = &c
	}

	// if all digits are the same, the CPF is not valid
	if allEq {
		return false
	}

	return checksum(ds[:9]) == ds[9] && checksum(ds[:10]) == ds[10]
}

func checksum(ds []int64) int64 {
	var s int64
	for i, n := range ds {
		s += n * int64(len(ds)+1-i)
	}
	r := 11 - (s % 11)
	if r == 10 {
		return 0
	}
	return r
}
