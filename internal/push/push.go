package push

import "context"

type Provider interface {
	Send(ctx context.Context, deviceID string, payload any) error
}

type Mock struct {}

func (m Mock) Send(ctx context.Context, deviceID string, payload any) error {
	return nil
}
