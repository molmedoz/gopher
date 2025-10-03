package env

import "os"

// Provider defines the interface for environment variable access
type Provider interface {
	Getenv(key string) string
}

// DefaultProvider implements Provider using os.Getenv
type DefaultProvider struct{}

// Getenv returns the value of the environment variable named by the key
func (d *DefaultProvider) Getenv(key string) string {
	return os.Getenv(key)
}

// MockProvider implements Provider for testing
type MockProvider struct {
	env map[string]string
}

// NewMockProvider creates a new MockProvider with the given environment variables
func NewMockProvider(env map[string]string) *MockProvider {
	mockEnv := &MockProvider{
		env: make(map[string]string),
	}
	// Copy the provided environment variables
	for key, value := range env {
		mockEnv.env[key] = value
	}
	return mockEnv
}

// Getenv returns the value of the environment variable from the mock
func (m *MockProvider) Getenv(key string) string {
	if value, exists := m.env[key]; exists {
		return value
	}
	return ""
}

// Setenv sets an environment variable in the mock
func (m *MockProvider) Setenv(key, value string) {
	if m.env == nil {
		m.env = make(map[string]string)
	}
	m.env[key] = value
}

// Clear removes all environment variables from the mock
func (m *MockProvider) Clear() {
	m.env = make(map[string]string)
}
