package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestIsAdminToken(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		realAdminToken string
		want           bool
	}{
		{
			name:           "success",
			token:          "abc",
			realAdminToken: "abc",
			want:           true,
		},
		{
			name:           "failed",
			token:          "abc",
			realAdminToken: "abcd",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := New(&Config{
				AdminToken: tt.realAdminToken},
				zap.NewNop(),
			)

			got := auth.IsAdminToken(tt.token)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		wantErr error
	}{
		{
			name:    "multi",
			login:   "12345678abc",
			wantErr: nil,
		},
		{
			name:    "only digits",
			login:   "12345678",
			wantErr: nil,
		},
		{
			name:    "only letters",
			login:   "abcdefgh",
			wantErr: nil,
		},
		{
			name:    "upper case",
			login:   "ABCDEFGH",
			wantErr: nil,
		},
		{
			name:    "small length",
			login:   "1234567",
			wantErr: ErrShortLogin,
		},
		{
			name:    "special symbol",
			login:   "1234567%",
			wantErr: ErrInvalidLogin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := New(&Config{}, zap.NewNop())

			err := auth.ValidateLogin(tt.login)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "success",
			password: "12345Aa#",
			wantErr:  nil,
		},
		{
			name:     "short password",
			password: "1234Aa#",
			wantErr:  ErrShortPassword,
		},
		{
			name:     "miss uppercase letter",
			password: "abcdefg1#",
			wantErr:  ErrMissingUpper,
		},
		{
			name:     "miss lowercase letter",
			password: "ABCDEFG1#",
			wantErr:  ErrMissingLower,
		},
		{
			name:     "miss digit",
			password: "abcdefG#",
			wantErr:  ErrMissingDigit,
		},
		{
			name:     "miss special symbol",
			password: "12345Aa1",
			wantErr:  ErrMissingSpecial,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := New(&Config{}, zap.NewNop())

			err := auth.ValidatePassword(tt.password)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
