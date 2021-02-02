package usecase

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"time"
)

// AuthUseCase is the interface that wraps all business logic methods related to authentication.
type AuthUseCase interface {
	Login(ctx context.Context, loginInput AuthLoginInput) (*AuthTokenOutput, error)
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
