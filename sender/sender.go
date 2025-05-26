package sender

import (
	"beats_refactor/config"
	"context"
)

type SenderService struct{}

func NewSenderService(c []*config.SenderConfig, ctx context.Context) (*SenderService, error) {
	return nil, nil
}

func check(c []*config.SenderConfig) error {
	return nil
}
