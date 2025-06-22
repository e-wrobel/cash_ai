package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Equal(t *testing.T) {
	tests := []struct {
		name string
		t1   Transaction
		t2   Transaction
		want bool
	}{
		{
			name: "identical transactions",
			t1:   Transaction{"1", "user1", 100.0, "credit", 1234567890},
			t2:   Transaction{"1", "user1", 100.0, "credit", 1234567890},
			want: true,
		},
		{
			name: "different ID",
			t1:   Transaction{"1", "user1", 100.0, "credit", 1234567890},
			t2:   Transaction{"2", "user1", 100.0, "credit", 1234567890},
			want: false,
		},
		{
			name: "different UserID",
			t1:   Transaction{"1", "user1", 100.0, "credit", 1234567890},
			t2:   Transaction{"1", "user2", 100.0, "credit", 1234567890},
			want: false,
		},
		{
			name: "different Amount",
			t1:   Transaction{"1", "user1", 100.0, "credit", 1234567890},
			t2:   Transaction{"1", "user1", 200.0, "credit", 1234567890},
			want: false,
		},
		{
			name: "different Type",
			t1:   Transaction{"1", "user1", 100.0, "credit", 1234567890},
			t2:   Transaction{"1", "user1", 100.0, "debit", 1234567890},
			want: false,
		},
		{
			name: "different Timestamp",
			t1:   Transaction{"1", "user1", 100.0, "credit", 1234567890},
			t2:   Transaction{"1", "user1", 100.0, "credit", 9876543210},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t1.Equal(tt.t2)
			assert.Equalf(t, tt.want, got, "Equal() = %v, want %v", got, tt.want)
		})
	}
}
