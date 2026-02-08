package totp

import (
	"context"
	"errors"
)

type Verifier interface {
	Verify(ctx context.Context, tenantID, userID, code string) (bool, error)
}

type Disabled struct {}

func (Disabled) Verify(ctx context.Context, tenantID, userID, code string) (bool, error) {
	return false, errors.New("totp_verifier_disabled")
}
