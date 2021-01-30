package model

import "testing"

func TestFloat64ToMoney(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want Money
	}{
		{
			name: "zero",
			f:    0,
			want: 0,
		},
		{
			name: "positive 1.0",
			f:    1.0,
			want: 100,
		},
		{
			name: "negative 1.0",
			f:    -1.0,
			want: -100,
		},
		{
			name: "positive 1123452.95",
			f:    1123452.95,
			want: 112345295,
		},
		{
			name: "negative 1123452.95",
			f:    -1123452.95,
			want: -112345295,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64ToMoney(tt.f); got != tt.want {
				t.Errorf("Float64ToMoney() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Float64(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want float64
	}{
		{
			name: "zero",
			m:    0,
			want: 0,
		},
		{
			name: "positive 100",
			m:    100,
			want: 1.0,
		},
		{
			name: "negative 100",
			m:    -100,
			want: -1.0,
		},
		{
			name: "positive 112345295",
			m:    112345295,
			want: 1123452.95,
		},
		{
			name: "negative 112345295",
			m:    -112345295,
			want: -1123452.95,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Float64(); got != tt.want {
				t.Errorf("Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Int64(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want int64
	}{
		{
			name: "zero",
			m:    0,
			want: 0,
		},
		{
			name: "positive 1",
			m:    1,
			want: 1,
		},
		{
			name: "negative 1",
			m:    -1,
			want: -1,
		},
		{
			name: "positive 112345295",
			m:    112345295,
			want: 112345295,
		},
		{
			name: "negative 112345295",
			m:    -112345295,
			want: -112345295,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Int64(); got != tt.want {
				t.Errorf("Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}
