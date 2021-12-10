package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

// AuthUseCase is the interface that wraps all business logic methods related to authentication.
type AuthUseCase interface {
	Login(ctx context.Context, loginInput AuthLoginInput) (*AuthTokenOutput, error)
	Authorize(ctx context.Context, accessToken string) (*jwt.StandardClaims, error)
}

type authUseCase struct {
	secretKey      string
	accessTokenDur time.Duration
	accRepo        repository.AccountRepository
}

// NewAuthUseCase instantiates a new AuthUseCase.
func NewAuthUseCase(
	secretKey string,
	accessTokenDur time.Duration,
	accRepo repository.AccountRepository,
) AuthUseCase {
	return &authUseCase{
		secretKey:      secretKey,
		accessTokenDur: accessTokenDur,
		accRepo:        accRepo,
	}
}
