package utils

import (
	"testing"
)

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Basic cases
		{"", ""},
		{"user", "User"},
		{"User", "User"},
		{"USER", "User"},

		// Snake case
		{"user_profile", "UserProfile"},
		{"user_profile_settings", "UserProfileSettings"},
		{"_private", "Private"},

		// Kebab case
		{"user-profile", "UserProfile"},
		{"user-profile-settings", "UserProfileSettings"},

		// Camel case
		{"userProfile", "UserProfile"},
		{"userProfileSettings", "UserProfileSettings"},

		// Mixed
		{"user_profile-settings", "UserProfileSettings"},
		{"USER_PROFILE", "UserProfile"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToPascalCase(tt.input)
			if got != tt.want {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Basic cases
		{"", ""},
		{"user", "user"},
		{"User", "user"},
		{"USER", "user"},

		// Snake case
		{"user_profile", "userProfile"},
		{"user_profile_settings", "userProfileSettings"},

		// Kebab case
		{"user-profile", "userProfile"},
		{"user-profile-settings", "userProfileSettings"},

		// Already camel case
		{"userProfile", "userProfile"},
		{"userProfileSettings", "userProfileSettings"},

		// Pascal case
		{"UserProfile", "userProfile"},
		{"UserProfileSettings", "userProfileSettings"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToCamelCase(tt.input)
			if got != tt.want {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Basic cases
		{"", ""},
		{"user", "user"},
		{"User", "user"},

		// Pascal case
		{"UserProfile", "user_profile"},
		{"UserProfileSettings", "user_profile_settings"},

		// Camel case
		{"userProfile", "user_profile"},
		{"userProfileSettings", "user_profile_settings"},

		// Already snake case
		{"user_profile", "user_profile"},

		// Kebab case
		{"user-profile", "user_profile"},

		// Acronyms
		{"userID", "user_id"},
		{"XMLParser", "xml_parser"},
		{"parseJSON", "parse_json"},

		// Multiple underscores (should be cleaned)
		{"user__profile", "user_profile"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToSnakeCase(tt.input)
			if got != tt.want {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Basic cases
		{"", ""},
		{"user", "user"},
		{"User", "user"},

		// Pascal case
		{"UserProfile", "user-profile"},
		{"UserProfileSettings", "user-profile-settings"},

		// Camel case
		{"userProfile", "user-profile"},
		{"userProfileSettings", "user-profile-settings"},

		// Snake case
		{"user_profile", "user-profile"},

		// Already kebab case
		{"user-profile", "user-profile"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToKebabCase(tt.input)
			if got != tt.want {
				t.Errorf("ToKebabCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToPackageName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"user", "user"},
		{"User", "user"},
		{"UserProfile", "userprofile"},
		{"user_profile", "userprofile"},
		{"user-profile", "userprofile"},
		{"User123", "user123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToPackageName(tt.input)
			if got != tt.want {
				t.Errorf("ToPackageName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToVariableName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"user", "user"},
		{"User", "user"},
		{"user_profile", "userProfile"},
		{"UserProfile", "userProfile"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToVariableName(tt.input)
			if got != tt.want {
				t.Errorf("ToVariableName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToRepoVariableName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user", "userRepo"},
		{"User", "userRepo"},
		{"userProfile", "userProfileRepo"},
		{"user_profile", "userProfileRepo"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToRepoVariableName(tt.input)
			if got != tt.want {
				t.Errorf("ToRepoVariableName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToServiceVariableName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user", "userService"},
		{"User", "userService"},
		{"userProfile", "userProfileService"},
		{"user_profile", "userProfileService"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToServiceVariableName(tt.input)
			if got != tt.want {
				t.Errorf("ToServiceVariableName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToControllerVariableName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user", "userController"},
		{"User", "userController"},
		{"userProfile", "userProfileController"},
		{"user_profile", "userProfileController"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToControllerVariableName(tt.input)
			if got != tt.want {
				t.Errorf("ToControllerVariableName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToModelName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user", "User"},
		{"User", "User"},
		{"user_profile", "UserProfile"},
		{"userProfile", "UserProfile"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToModelName(tt.input)
			if got != tt.want {
				t.Errorf("ToModelName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToTableName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user", "users"},
		{"User", "users"},
		{"person", "people"},
		{"category", "categories"},
		{"userProfile", "user_profiles"},
		{"order_item", "order_items"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToTableName(tt.input)
			if got != tt.want {
				t.Errorf("ToTableName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToURLPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user", "/users"},
		{"User", "/users"},
		{"person", "/people"},
		{"category", "/categories"},
		{"userProfile", "/user-profiles"},
		{"order_item", "/order-items"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToURLPath(tt.input)
			if got != tt.want {
				t.Errorf("ToURLPath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"user", "users"},
		{"category", "categories"},
		{"person", "people"},
		{"child", "children"},
		{"status", "statuses"},
		{"box", "boxes"},
		{"city", "cities"},
		{"country", "countries"},
		{"leaf", "leafs"}, // inflection library doesn't handle irregular plurals like leaf->leaves
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Pluralize(tt.input)
			if got != tt.want {
				t.Errorf("Pluralize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSingularize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"users", "user"},
		{"categories", "category"},
		{"people", "person"},
		{"children", "child"},
		{"statuses", "status"},
		{"boxes", "box"},
		{"cities", "city"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Singularize(tt.input)
			if got != tt.want {
				t.Errorf("Singularize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToJSONTag(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Name", "name"},
		{"FirstName", "first_name"},
		{"userID", "user_id"},
		{"user_name", "user_name"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToJSONTag(tt.input)
			if got != tt.want {
				t.Errorf("ToJSONTag(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToLabel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"name", "Name"},
		{"firstName", "First Name"},
		{"first_name", "First Name"},
		{"userID", "User ID"},
		{"parseJSON", "Parse JSON"},
		{"XMLParser", "XML Parser"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToLabel(tt.input)
			if got != tt.want {
				t.Errorf("ToLabel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"validName", "validName"},
		{"Valid_Name", "Valid_Name"},
		{"123name", "_23name"},
		{"my-name", "myname"},
		{"my name", "myname"},
		{"my.name", "myname"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := SanitizeIdentifier(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeIdentifier(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
