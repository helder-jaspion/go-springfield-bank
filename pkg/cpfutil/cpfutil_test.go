package cpfutil

import "testing"

func TestClean(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cpf  string
		want string
	}{
		{
			name: "valid masked CPF",
			cpf:  "854.725.670-92",
			want: "85472567092",
		},
		{
			name: "valid unmasked CPF",
			cpf:  "94640164009",
			want: "94640164009",
		},
		{
			name: "incomplete masked CPF",
			cpf:  "854.725.670-9",
			want: "8547256709",
		},
		{
			name: "incomplete unmasked CPF",
			cpf:  "8547256709",
			want: "8547256709",
		},
		{
			name: "too much digits masked CPF",
			cpf:  "854.725.670-999",
			want: "854725670999",
		},
		{
			name: "too much digits unmasked CPF",
			cpf:  "854725670999",
			want: "854725670999",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Clean(tt.cpf); got != tt.want {
				t.Errorf("Clean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cpf  string
		want string
	}{
		{
			name: "valid masked CPF",
			cpf:  "854.725.670-92",
			want: "854.725.670-92",
		},
		{
			name: "valid unmasked CPF",
			cpf:  "94640164009",
			want: "946.401.640-09",
		},
		{
			name: "incomplete masked CPF, return input",
			cpf:  "854.725.670-9",
			want: "854.725.670-9",
		},
		{
			name: "incomplete unmasked CPF, return input",
			cpf:  "8547256709",
			want: "8547256709",
		},
		{
			name: "too much digits masked CPF, return input",
			cpf:  "854.725.670-999",
			want: "854.725.670-999",
		},
		{
			name: "too much digits unmasked CPF, return input",
			cpf:  "854725670999",
			want: "854725670999",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Format(tt.cpf); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cpf  string
		want bool
	}{
		{
			name: "masked CPF - valid",
			cpf:  "854.725.670-92",
			want: true,
		},
		{
			name: "unmasked CPF - valid",
			cpf:  "94640164009",
			want: true,
		},
		{
			name: "incomplete masked CPF - invalid",
			cpf:  "854.725.670-9",
			want: false,
		},
		{
			name: "incomplete unmasked CPF - invalid",
			cpf:  "8547256709",
			want: false,
		},
		{
			name: "all digits equal CPF - invalid",
			cpf:  "11111111111",
			want: false,
		},
		{
			name: "invalid masked CPF - invalid",
			cpf:  "123.123.123-12",
			want: false,
		},
		{
			name: "invalid unmasked CPF - invalid",
			cpf:  "12312312312",
			want: false,
		},
		{
			name: "empty CPF - invalid",
			cpf:  "",
			want: false,
		},
		{
			name: "too much digits CPF - invalid",
			cpf:  "854.725.670-921",
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := IsValid(tt.cpf); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
