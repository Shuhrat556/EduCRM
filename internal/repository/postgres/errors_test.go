package postgres

import (
	"errors"
	"testing"
)

func TestIsUniqueViolation_stringFallback(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{errors.New("duplicate key value violates unique constraint"), true},
		{errors.New("UNIQUE constraint failed"), true},
		{errors.New("other error"), false},
	}
	for _, tt := range tests {
		if got := isUniqueViolation(tt.err); got != tt.want {
			t.Errorf("%q: got %v want %v", tt.err.Error(), got, tt.want)
		}
	}
}

func TestIsForeignKeyViolation_stringFallback(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{errors.New("violates foreign key constraint"), true},
		{errors.New("foreign key violation"), true},
		{errors.New("nope"), false},
	}
	for _, tt := range tests {
		if got := isForeignKeyViolation(tt.err); got != tt.want {
			t.Errorf("%q: got %v want %v", tt.err.Error(), got, tt.want)
		}
	}
}
