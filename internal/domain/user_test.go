package domain

import "testing"

func TestNewUser(t *testing.T) {
	_, err := NewUser("", "Name", "email@example.com", "password")
	if err == nil {
		t.Error("expected error for empty ID")
	}

	user, err := NewUser("1", "Name", "email@example.com", "password")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if user.ID != "1" {
		t.Errorf("expected ID '1', ogt %s", user.ID)
	}
}
