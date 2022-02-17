package model

import "testing"

func TestNewCPF(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    string
		want CPF
	}{
		{
			name: "masked CPF",
			c:    "854.725.670-92",
			want: CPF("85472567092"),
		},
		{
			name: "unmasked CPF",
			c:    "94640164009",
			want: CPF("94640164009"),
		},
		{
			name: "incomplete masked CPF",
			c:    "854.725.670-9",
			want: CPF("8547256709"),
		},
		{
			name: "incomplete unmasked CPF",
			c:    "8547256709",
			want: CPF("8547256709"),
		},
		{
			name: "empty CPF",
			c:    "",
			want: CPF(""),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewCPF(tt.c); got != tt.want {
				t.Errorf("NewCPF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCPF_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    CPF
		want string
	}{
		{
			name: "valid masked CPF",
			c:    "854.725.670-92",
			want: "854.725.670-92",
		},
		{
			name: "valid unmasked CPF",
			c:    "94640164009",
			want: "946.401.640-09",
		},
		{
			name: "incomplete masked CPF, return input",
			c:    "854.725.670-9",
			want: "854.725.670-9",
		},
		{
			name: "incomplete unmasked CPF, return input",
			c:    "8547256709",
			want: "8547256709",
		},
		{
			name: "too much digits masked CPF, return input",
			c:    "854.725.670-999",
			want: "854.725.670-999",
		},
		{
			name: "too much digits unmasked CPF, return input",
			c:    "854725670999",
			want: "854725670999",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
