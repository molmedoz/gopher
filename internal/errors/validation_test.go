package errors

import (
	"testing"
)

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{"valid version", "1.21.0", false},
		{"valid version with go prefix", "go1.21.0", false},
		{"valid version with beta", "1.21.0-beta1", false},
		{"valid version with rc", "1.21.0-rc1", false},
		{"empty version", "", true},
		{"invalid format", "invalid", true},
		{"too many parts", "1.2.3.4", true},
		{"zero version", "0.0.0", true},
		{"missing minor", "1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAliasName(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		wantErr bool
	}{
		{"valid alias", "stable", false},
		{"valid with numbers", "go1.21", false},
		{"valid with hyphens", "my-alias", false},
		{"valid with underscores", "my_alias", false},
		{"valid with dots", "my.alias", false},
		{"empty alias", "", true},
		{"too long", "this-alias-name-is-way-too-long-and-exceeds-the-fifty-character-limit", true},
		{"reserved name", "system", true},
		{"reserved name case insensitive", "SYSTEM", true},
		{"invalid characters", "my@alias", true},
		{"starts with hyphen", "-alias", true},
		{"starts with underscore", "_alias", true},
		{"starts with dot", ".alias", true},
		{"ends with hyphen", "alias-", true},
		{"ends with underscore", "alias_", true},
		{"ends with dot", "alias.", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAliasName(tt.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAliasName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateConfigValue(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{"valid gopath_mode", "gopath_mode", "shared", false},
		{"valid gopath_mode version-specific", "gopath_mode", "version-specific", false},
		{"valid gopath_mode custom", "gopath_mode", "custom", false},
		{"invalid gopath_mode", "gopath_mode", "invalid", true},
		{"valid set_environment true", "set_environment", "true", false},
		{"valid set_environment false", "set_environment", "false", false},
		{"invalid set_environment", "set_environment", "yes", true},
		{"valid auto_cleanup true", "auto_cleanup", "true", false},
		{"valid auto_cleanup false", "auto_cleanup", "false", false},
		{"invalid auto_cleanup", "auto_cleanup", "yes", true},
		{"valid max_versions", "max_versions", "5", false},
		{"empty max_versions", "max_versions", "", true},
		{"valid mirror_url", "mirror_url", "https://golang.org/dl/", false},
		{"valid mirror_url http", "mirror_url", "http://example.com", false},
		{"invalid mirror_url", "mirror_url", "ftp://example.com", true},
		{"empty mirror_url", "mirror_url", "", true},
		{"valid custom_gopath", "custom_gopath", "/path/to/gopath", false},
		{"empty custom_gopath", "custom_gopath", "", true},
		{"unknown config option", "unknown_option", "value", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfigValue(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid path", "/path/to/file", false},
		{"valid windows path", "C:\\path\\to\\file", false},
		{"empty path", "", true},
		{"path with <", "/path<file", true},
		{"path with >", "/path>file", true},
		{"path with :", "/path:file", true},
		{"path with \"", "/path\"file", true},
		{"path with |", "/path|file", true},
		{"path with ?", "/path?file", true},
		{"path with *", "/path*file", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateShell(t *testing.T) {
	tests := []struct {
		name    string
		shell   string
		wantErr bool
	}{
		{"valid bash", "bash", false},
		{"valid zsh", "zsh", false},
		{"valid fish", "fish", false},
		{"valid powershell", "powershell", false},
		{"valid cmd", "cmd", false},
		{"invalid shell", "invalid", true},
		{"empty shell", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateShell(tt.shell)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateShell() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateKeyValuePair(t *testing.T) {
	tests := []struct {
		name     string
		keyValue string
		wantErr  bool
	}{
		{"valid pair", "key=value", false},
		{"valid with equals in value", "key=value=with=equals", false},
		{"empty pair", "", true},
		{"no equals", "keyvalue", true},
		{"empty key", "=value", true},
		{"empty value", "key=", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateKeyValuePair(tt.keyValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateKeyValuePair() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
		wantErr bool
	}{
		{"valid install with args", "install", []string{"1.21.0"}, false},
		{"valid uninstall with args", "uninstall", []string{"1.21.0"}, false},
		{"valid use with args", "use", []string{"1.21.0"}, false},
		{"valid show with args", "show", []string{"1.21.0"}, false},
		{"valid set with args", "set", []string{"key=value"}, false},
		{"valid alias with args", "alias", []string{"list"}, false},
		{"valid env with args", "env", []string{"show"}, false},
		{"install without args", "install", []string{}, true},
		{"uninstall without args", "uninstall", []string{}, true},
		{"use without args", "use", []string{}, true},
		{"show without args", "show", []string{}, true},
		{"set without args", "set", []string{}, true},
		{"alias without args", "alias", []string{}, true},
		{"env without args", "env", []string{}, true},
		{"empty command", "", []string{}, true},
		{"valid list without args", "list", []string{}, false},
		{"valid help without args", "help", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommand(tt.command, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
