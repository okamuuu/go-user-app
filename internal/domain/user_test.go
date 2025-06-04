package domain

import "testing"

func TestNewUser(t *testing.T) {
	user, err := NewUser("Name", "email@example.com", "password")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if user.ID != 0 {
		t.Errorf("expected ID 1, ogt %v", user.ID)
	}
}
