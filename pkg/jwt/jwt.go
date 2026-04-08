package jwt

import (
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

const tokenTypeAccess = "access"

// AccessClaims is the JWT payload for API access tokens.
type AccessClaims struct {
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwtlib.RegisteredClaims
}

// Manager signs and parses access JWTs.
type Manager struct {
	secret     []byte
	expiration time.Duration
	issuer     string
}

// NewManager creates a JWT manager for access tokens.
func NewManager(secret string, accessExpiration time.Duration, issuer string) *Manager {
	return &Manager{
		secret:     []byte(secret),
		expiration: accessExpiration,
		issuer:     issuer,
	}
}

// GenerateAccessToken creates a signed access JWT with subject (user ID) and role.
func (m *Manager) GenerateAccessToken(userID string, role string) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		Role:      role,
		TokenType: tokenTypeAccess,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(m.expiration)),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseAccessToken validates an access JWT and returns claims.
func (m *Manager) ParseAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &AccessClaims{}, func(t *jwtlib.Token) (any, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims.TokenType != tokenTypeAccess {
		return nil, fmt.Errorf("invalid token type")
	}
	return claims, nil
}

// AccessTTL returns the configured access token lifetime.
func (m *Manager) AccessTTL() time.Duration {
	return m.expiration
}
