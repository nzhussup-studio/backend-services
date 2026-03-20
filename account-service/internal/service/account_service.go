package service

import (
	"context"

	"account-service/internal/keycloak"
)

type adminClient interface {
	LogoutUser(ctx context.Context, userID string) error
	DeleteUser(ctx context.Context, userID string) error
	ResolveUserID(ctx context.Context, subject string, username string, email string) (string, error)
}

type AccountService struct {
	adminClient adminClient
}

func NewAccountService(adminClient *keycloak.AdminClient) *AccountService {
	return &AccountService{adminClient: adminClient}
}

func (s *AccountService) DeleteCurrentAccount(ctx context.Context, subject string, username string, email string) error {
	userID, err := s.adminClient.ResolveUserID(ctx, subject, username, email)
	if err != nil {
		return err
	}

	if err := s.adminClient.LogoutUser(ctx, userID); err != nil {
		return err
	}

	if err := s.adminClient.DeleteUser(ctx, userID); err != nil {
		return err
	}

	return nil
}
