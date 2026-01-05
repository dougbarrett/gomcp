package utils

import (
	"strings"
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		// Valid cases
		{"valid simple", "myproject", false, ""},
		{"valid with hyphen", "my-project", false, ""},
		{"valid with underscore", "my_project", false, ""},
		{"valid with numbers", "project123", false, ""},
		{"valid mixed", "my-project_123", false, ""},
		{"valid uppercase", "MyProject", false, ""},

		// Invalid cases
		{"empty string", "", true, "project name is required"},
		{"starts with number", "123project", true, "must start with a letter"},
		{"starts with hyphen", "-project", true, "must start with a letter"},
		{"starts with underscore", "_project", true, "must start with a letter"},
		{"special chars at sign", "my@project", true, "must start with a letter"},
		{"special chars dollar", "my$project", true, "must start with a letter"},
		{"spaces", "my project", true, "must start with a letter"},
		{"dots", "my.project", true, "must start with a letter"},
		{"reserved word break", "break", true, "reserved word"},
		{"reserved word func", "func", true, "reserved word"},
		{"reserved word package", "package", true, "reserved word"},
		{"too long", strings.Repeat("a", 200), true, "too long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateProjectName(%q) error = %v, want error containing %q", tt.input, err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateModulePath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		// Valid cases
		{"simple", "mymodule", false, ""},
		{"with domain", "github.com/user/project", false, ""},
		{"nested path", "github.com/user/project/subdir", false, ""},
		{"with numbers", "github.com/user123/project456", false, ""},
		{"with hyphens", "github.com/my-user/my-project", false, ""},
		{"with underscores", "github.com/my_user/my_project", false, ""},
		{"short domain", "example.com/pkg", false, ""},

		// Invalid cases
		{"empty string", "", true, "module path is required"},
		{"with spaces", "github.com/user name/project", true, "whitespace"},
		{"with tabs", "github.com/user\tname/project", true, "whitespace"},
		{"double slashes", "github.com//project", true, "invalid module path"},
		{"trailing slash", "github.com/user/project/", true, "invalid module path"},
		{"too long", "github.com/" + strings.Repeat("a", 300), true, "too long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModulePath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateModulePath(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateModulePath(%q) error = %v, want error containing %q", tt.input, err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateDatabaseType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty (default)", "", false},
		{"sqlite", "sqlite", false},
		{"postgres", "postgres", false},
		{"mysql", "mysql", false},
		{"invalid mongodb", "mongodb", true},
		{"invalid oracle", "oracle", true},
		{"invalid uppercase", "SQLITE", true},
		{"invalid mixed case", "Postgres", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDatabaseType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDatabaseType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDomainName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		// Valid cases
		{"simple lowercase", "user", false, ""},
		{"with underscore", "user_profile", false, ""},
		{"with numbers", "order123", false, ""},
		{"camelCase", "userProfile", false, ""},
		{"PascalCase", "UserProfile", false, ""},
		{"single char", "x", false, ""},
		{"starts with underscore", "_private", false, ""},

		// Invalid cases
		{"empty string", "", true, "domain name is required"},
		{"starts with number", "123user", true, "must start with a letter"},
		{"with hyphen", "user-profile", true, "invalid character"},
		{"with space", "user profile", true, "invalid character"},
		{"with dot", "user.profile", true, "invalid character"},
		{"reserved word if", "if", true, "reserved word"},
		{"reserved word for", "for", true, "reserved word"},
		{"reserved word internal", "internal", true, "reserved name"},
		{"reserved word vendor", "vendor", true, "reserved name"},
		{"too long", strings.Repeat("a", 100), true, "too long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDomainName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDomainName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateDomainName(%q) error = %v, want error containing %q", tt.input, err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateFieldName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		// Valid cases
		{"simple", "Name", false, ""},
		{"with underscore", "First_Name", false, ""},
		{"with numbers", "Field123", false, ""},
		{"single letter", "X", false, ""},
		{"ID style", "ID", false, ""},
		{"URL style", "URL", false, ""},

		// Invalid cases
		{"empty string", "", true, "field name is required"},
		{"lowercase start", "name", true, "uppercase letter"},
		{"starts with underscore", "_Name", true, "uppercase letter"},
		{"starts with number", "123Name", true, "not a valid Go identifier"},
		{"with hyphen", "First-Name", true, "not a valid Go identifier"},
		{"with space", "First Name", true, "not a valid Go identifier"},
		{"reserved word", "Type", true, "reserved word"},
		{"too long", strings.Repeat("A", 100), true, "too long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFieldName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFieldName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateFieldName(%q) error = %v, want error containing %q", tt.input, err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateFieldType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid built-in types
		{"string", "string", false},
		{"int", "int", false},
		{"int64", "int64", false},
		{"uint", "uint", false},
		{"float64", "float64", false},
		{"bool", "bool", false},
		{"time.Time", "time.Time", false},
		{"*time.Time", "*time.Time", false},
		{"[]byte", "[]byte", false},
		{"[]string", "[]string", false},

		// Valid custom types
		{"custom type", "CustomType", false},
		{"package.Type", "models.Status", false},
		{"pointer to custom", "*CustomType", false},
		{"slice of custom", "[]CustomType", false},

		// Invalid cases
		{"empty string", "", true},
		{"invalid chars", "my-type", true},
		{"spaces", "my type", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFieldType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFieldType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateFormType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid types
		{"empty (default)", "", false},
		{"input", "input", false},
		{"textarea", "textarea", false},
		{"select", "select", false},
		{"checkbox", "checkbox", false},
		{"switch", "switch", false},
		{"date", "date", false},
		{"time", "time", false},
		{"datetime", "datetime", false},
		{"email", "email", false},
		{"password", "password", false},
		{"number", "number", false},
		{"rating", "rating", false},
		{"tags", "tags", false},
		{"slider", "slider", false},

		// Invalid types
		{"invalid radio", "radio", true},
		{"invalid file", "file", true},
		{"invalid hidden", "hidden", true},
		{"invalid uppercase", "INPUT", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFormType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFormType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateViewType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid types
		{"list", "list", false},
		{"show", "show", false},
		{"form", "form", false},
		{"card", "card", false},
		{"table", "table", false},
		{"custom", "custom", false},

		// Invalid types
		{"empty", "", true},
		{"invalid grid", "grid", true},
		{"invalid detail", "detail", true},
		{"invalid uppercase", "LIST", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateViewType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateViewType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateModalType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid types
		{"dialog", "dialog", false},
		{"sheet", "sheet", false},
		{"confirm", "confirm", false},

		// Invalid types
		{"empty", "", true},
		{"invalid popup", "popup", true},
		{"invalid modal", "modal", true},
		{"invalid alert", "alert", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModalType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateModalType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateConfigType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid types
		{"page", "page", false},
		{"menu", "menu", false},
		{"app", "app", false},
		{"messages", "messages", false},

		// Invalid types
		{"empty", "", true},
		{"invalid config", "config", true},
		{"invalid settings", "settings", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfigType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateLayoutType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid types
		{"empty (default)", "", false},
		{"default", "default", false},
		{"dashboard", "dashboard", false},
		{"landing", "landing", false},
		{"blank", "blank", false},

		// Invalid types
		{"invalid admin", "admin", true},
		{"invalid full", "full", true},
		{"invalid uppercase", "DEFAULT", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLayoutType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLayoutType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateHTTPMethod(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid methods (case insensitive)
		{"GET uppercase", "GET", false},
		{"POST uppercase", "POST", false},
		{"PUT uppercase", "PUT", false},
		{"PATCH uppercase", "PATCH", false},
		{"DELETE uppercase", "DELETE", false},
		{"get lowercase", "get", false},
		{"post lowercase", "post", false},
		{"mixed case", "Get", false},

		// Invalid methods
		{"empty", "", true},
		{"HEAD", "HEAD", true},
		{"OPTIONS", "OPTIONS", true},
		{"CONNECT", "CONNECT", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHTTPMethod(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHTTPMethod(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateURLPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid paths
		{"empty (optional)", "", false},
		{"root", "/", false},
		{"simple path", "/users", false},
		{"nested path", "/api/v1/users", false},
		{"with param", "/users/{id}", false},
		{"with multiple params", "/users/{id}/posts/{postId}", false},
		{"with hyphen", "/user-profiles", false},
		{"with underscore", "/user_profiles", false},

		// Invalid paths
		{"no leading slash", "users", true},
		{"with spaces", "/my users", true},
		{"with at sign", "/users/@me", true},
		{"with query", "/users?id=1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURLPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURLPath(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateLocale(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid locales
		{"empty (default)", "", false},
		{"en", "en", false},
		{"fr", "fr", false},
		{"de", "de", false},
		{"en-US", "en-US", false},
		{"pt-BR", "pt-BR", false},
		{"zh-CN", "zh-CN", false},

		// Invalid locales
		{"uppercase", "EN", true},
		{"three letters", "eng", true},
		{"lowercase region", "en-us", true},
		{"underscore", "en_US", true},
		{"single letter", "e", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLocale(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLocale(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateComponentName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid names
		{"simple", "Button", false},
		{"lowercase", "button", false},
		{"with underscore", "my_button", false},
		{"with numbers", "Button123", false},
		{"camelCase", "myButton", false},
		{"PascalCase", "MyButton", false},

		// Invalid names
		{"empty", "", true},
		{"with hyphen", "my-button", true},
		{"with space", "my button", true},
		{"starts with number", "123Button", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateComponentName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateComponentName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateRelationshipType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		// Valid types
		{"belongs_to", "belongs_to", false, ""},
		{"has_one", "has_one", false, ""},
		{"has_many", "has_many", false, ""},
		{"many_to_many", "many_to_many", false, ""},

		// Invalid types
		{"empty", "", true, "relationship type is required"},
		{"invalid type", "invalid", true, "invalid relationship type"},
		{"typo belongs", "belongsto", true, "invalid relationship type"},
		{"uppercase", "BELONGS_TO", true, "invalid relationship type"},
		{"camelCase", "belongsTo", true, "invalid relationship type"},
		{"one_to_one", "one_to_one", true, "invalid relationship type"},
		{"one_to_many", "one_to_many", true, "invalid relationship type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRelationshipType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRelationshipType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateRelationshipType(%q) error = %v, want error containing %q", tt.input, err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateRelationshipModel(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		// Valid model names (PascalCase)
		{"simple", "User", false, ""},
		{"two words", "UserProfile", false, ""},
		{"with numbers", "User123", false, ""},
		{"acronym", "APIKey", false, ""},
		{"longer name", "OrderLineItem", false, ""},

		// Invalid model names
		{"empty", "", true, "related model name is required"},
		{"lowercase", "user", true, "must be in PascalCase"},
		{"camelCase", "userProfile", true, "must be in PascalCase"},
		{"snake_case", "user_profile", true, "must be in PascalCase"},
		{"kebab-case", "user-profile", true, "not a valid identifier"},
		{"with spaces", "User Profile", true, "not a valid identifier"},
		{"starts with number", "123User", true, "not a valid identifier"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRelationshipModel(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRelationshipModel(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateRelationshipModel(%q) error = %v, want error containing %q", tt.input, err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateOnDelete(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		// Valid actions
		{"empty (defaults to CASCADE)", "", false, ""},
		{"CASCADE", "CASCADE", false, ""},
		{"SET NULL", "SET NULL", false, ""},
		{"RESTRICT", "RESTRICT", false, ""},
		{"NO ACTION", "NO ACTION", false, ""},
		{"lowercase cascade", "cascade", false, ""},
		{"lowercase set null", "set null", false, ""},
		{"mixed case", "Cascade", false, ""},

		// Invalid actions
		{"invalid action", "DELETE", true, "invalid ON DELETE action"},
		{"typo", "CASCAD", true, "invalid ON DELETE action"},
		{"SET_NULL with underscore", "SET_NULL", true, "invalid ON DELETE action"},
		{"random text", "something", true, "invalid ON DELETE action"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOnDelete(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOnDelete(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateOnDelete(%q) error = %v, want error containing %q", tt.input, err, tt.errMsg)
				}
			}
		})
	}
}
