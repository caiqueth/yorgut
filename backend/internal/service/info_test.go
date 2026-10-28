package service

import "testing"

// The administrator's e-mail comes from ADMIN_EMAIL, and is null without it.
func TestInstanceInfo(t *testing.T) {
	t.Setenv("ADMIN_EMAIL", " admin@test.com ")
	if got := NewInfoService().FindInstanceInfo().AdminEmail; got == nil || *got != "admin@test.com" {
		t.Errorf("AdminEmail = %v, want admin@test.com", got)
	}

	t.Setenv("ADMIN_EMAIL", "")
	if got := NewInfoService().FindInstanceInfo().AdminEmail; got != nil {
		t.Errorf("AdminEmail = %q, want nil", *got)
	}
}
