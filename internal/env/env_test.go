package env

import (
	"testing"
)

func TestDefaultProvider(t *testing.T) {
	provider := &DefaultProvider{}

	// Test that it returns actual environment variables
	path := provider.Getenv("PATH")
	if path == "" {
		t.Error("Expected PATH environment variable to be set")
	}

	// Test that it returns empty string for non-existent variables
	nonExistent := provider.Getenv("NON_EXISTENT_VAR_12345")
	if nonExistent != "" {
		t.Errorf("Expected empty string for non-existent variable, got: %s", nonExistent)
	}
}

func TestMockProvider(t *testing.T) {
	// Test with initial environment variables
	env := map[string]string{
		"TEST_VAR": "test_value",
		"PATH":     "/usr/local/bin:/usr/bin",
	}
	provider := NewMockProvider(env)

	// Test getting existing variables
	if value := provider.Getenv("TEST_VAR"); value != "test_value" {
		t.Errorf("Expected 'test_value', got: %s", value)
	}

	if value := provider.Getenv("PATH"); value != "/usr/local/bin:/usr/bin" {
		t.Errorf("Expected '/usr/local/bin:/usr/bin', got: %s", value)
	}

	// Test getting non-existent variable
	if value := provider.Getenv("NON_EXISTENT"); value != "" {
		t.Errorf("Expected empty string for non-existent variable, got: %s", value)
	}
}

func TestMockProvider_Setenv(t *testing.T) {
	provider := NewMockProvider(nil)

	// Test setting a new variable
	provider.Setenv("NEW_VAR", "new_value")
	if value := provider.Getenv("NEW_VAR"); value != "new_value" {
		t.Errorf("Expected 'new_value', got: %s", value)
	}

	// Test updating an existing variable
	provider.Setenv("NEW_VAR", "updated_value")
	if value := provider.Getenv("NEW_VAR"); value != "updated_value" {
		t.Errorf("Expected 'updated_value', got: %s", value)
	}
}

func TestMockProvider_Clear(t *testing.T) {
	provider := NewMockProvider(map[string]string{
		"VAR1": "value1",
		"VAR2": "value2",
	})

	// Verify variables exist
	if provider.Getenv("VAR1") != "value1" {
		t.Error("Expected VAR1 to be set")
	}

	// Clear and verify variables are gone
	provider.Clear()
	if provider.Getenv("VAR1") != "" {
		t.Error("Expected VAR1 to be cleared")
	}
	if provider.Getenv("VAR2") != "" {
		t.Error("Expected VAR2 to be cleared")
	}
}
