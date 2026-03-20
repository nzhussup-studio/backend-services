package service

import (
	"context"
	"errors"
	"testing"
)

type adminClientStub struct {
	logoutErr error
	deleteErr error
	loggedOut []string
	deleted   []string
}

func (s *adminClientStub) LogoutUser(_ context.Context, userID string) error {
	s.loggedOut = append(s.loggedOut, userID)
	return s.logoutErr
}

func (s *adminClientStub) DeleteUser(_ context.Context, userID string) error {
	s.deleted = append(s.deleted, userID)
	return s.deleteErr
}

func TestDeleteCurrentAccountLogsOutThenDeletes(t *testing.T) {
	stub := &adminClientStub{}
	service := &AccountService{adminClient: stub}

	err := service.DeleteCurrentAccount(context.Background(), "user-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(stub.loggedOut) != 1 || stub.loggedOut[0] != "user-123" {
		t.Fatalf("expected logout to be called for user-123, got %v", stub.loggedOut)
	}

	if len(stub.deleted) != 1 || stub.deleted[0] != "user-123" {
		t.Fatalf("expected delete to be called for user-123, got %v", stub.deleted)
	}
}

func TestDeleteCurrentAccountStopsWhenLogoutFails(t *testing.T) {
	stub := &adminClientStub{logoutErr: errors.New("logout failed")}
	service := &AccountService{adminClient: stub}

	err := service.DeleteCurrentAccount(context.Background(), "user-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if len(stub.deleted) != 0 {
		t.Fatalf("expected delete not to be called, got %v", stub.deleted)
	}
}
