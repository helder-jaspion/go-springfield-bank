package model

// Money represents monetary amount.
// It is an integer to prevent floating point math problems.
type Money int64

// Float64 converts Money to float64
func (m Money) Float64() float64 {
	return float64(m) / 100
}

// Int64 converts Money to int64
func (m Money) Int64() int64 {
	return int64(m)
}

// Float64ToMoney converts float64 to Money
func Float64ToMoney(f float64) Money {
	return Money(f * 100)
}
