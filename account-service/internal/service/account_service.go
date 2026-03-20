package service

import (
	"context"

	"account-service/internal/keycloak"
)

type adminClient interface {
	LogoutUser(ctx context.Context, userID string) error
	DeleteUser(ctx context.Context, userID string) error
}

type AccountService struct {
	adminClient adminClient
}

func NewAccountService(adminClient *keycloak.AdminClient) *AccountService {
	return &AccountService{adminClient: adminClient}
}

func (s *AccountService) DeleteCurrentAccount(ctx context.Context, subject string) error {
	if err := s.adminClient.LogoutUser(ctx, subject); err != nil {
		return err
	}

	if err := s.adminClient.DeleteUser(ctx, subject); err != nil {
		return err
	}

	return nil
}
