package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

var (
	// ErrAuthInvalidAccessToken happens when the JWT token is invalid or expired.
	ErrAuthInvalidAccessToken = errors.New("invalid access token")
)

// Authorize parses and verifies JWT token.
func (authUC authUseCase) Authorize(ctx context.Context, accessToken string) (*jwt.StandardClaims, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	token, err := jwt.ParseWithClaims(
		accessToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}

			return []byte(authUC.secretKey), nil
		},
	)
	if err != nil {
		log.Ctx(ctx).Err(err).Str("accessToken", accessToken).Msg("error parsing access token")
		return nil, ErrAuthInvalidAccessToken
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, ErrAuthInvalidAccessToken
	}

	return claims, nil
}
