package model

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTransfer(t *testing.T) {
	t.Parallel()

	type args struct {
		accountOriginID      string
		accountDestinationID string
		amount               float64
	}
	tests := []struct {
		name string
		args args
		want *Transfer
	}{
		{
			name: "positive amount success",
			args: args{
				accountOriginID:      "uuid-1",
				accountDestinationID: "uuid-2",
				amount:               10,
			},
			want: &Transfer{
				ID:                   "",
				AccountOriginID:      "uuid-1",
				AccountDestinationID: "uuid-2",
				Amount:               1000,
				CreatedAt:            time.Time{},
			},
		},
		{
			name: "negative amount success",
			args: args{
				accountOriginID:      "uuid-1",
				accountDestinationID: "uuid-2",
				amount:               -10.9,
			},
			want: &Transfer{
				ID:                   "",
				AccountOriginID:      "uuid-1",
				AccountDestinationID: "uuid-2",
				Amount:               -1090,
				CreatedAt:            time.Time{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewTransfer(tt.args.accountOriginID, tt.args.accountDestinationID, tt.args.amount)
			if len(got.ID) <= 0 {
				t.Errorf("NewTransfer() = %v, ID should not be empty", got)
			}
			got.ID = ""

			if got.CreatedAt.Before(time.Now().Add(-5 * time.Second)) {
				t.Errorf("NewTransfer() got = %v, want CreatedAt in the last 5 seconds", got)
			}
			got.CreatedAt = time.Time{}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransfer() = %v, want %v", got, tt.want)
			}
		})
	}
}
