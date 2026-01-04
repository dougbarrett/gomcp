package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
)

// ToPascalCase converts a string to PascalCase.
// Examples: "user_profile" -> "UserProfile", "user-profile" -> "UserProfile"
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// Split on common separators
	words := splitWords(s)
	var result strings.Builder

	for _, word := range words {
		if word == "" {
			continue
		}
		// Capitalize first letter, lowercase rest
		runes := []rune(strings.ToLower(word))
		runes[0] = unicode.ToUpper(runes[0])
		result.WriteString(string(runes))
	}

	return result.String()
}

// ToCamelCase converts a string to camelCase.
// Examples: "user_profile" -> "userProfile", "UserProfile" -> "userProfile"
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if pascal == "" {
		return ""
	}
	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// ToSnakeCase converts a string to snake_case.
// Examples: "UserProfile" -> "user_profile", "userProfile" -> "user_profile"
func ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				// Check if previous char was lowercase or if next char is lowercase
				prevLower := unicode.IsLower(runes[i-1])
				nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
				if prevLower || nextLower {
					result.WriteRune('_')
				}
			}
			result.WriteRune(unicode.ToLower(r))
		} else if r == '-' || r == ' ' {
			result.WriteRune('_')
		} else {
			result.WriteRune(r)
		}
	}

	// Clean up multiple underscores
	cleaned := regexp.MustCompile(`_+`).ReplaceAllString(result.String(), "_")
	return strings.Trim(cleaned, "_")
}

// ToKebabCase converts a string to kebab-case.
// Examples: "UserProfile" -> "user-profile", "user_profile" -> "user-profile"
func ToKebabCase(s string) string {
	snake := ToSnakeCase(s)
	return strings.ReplaceAll(snake, "_", "-")
}

// ToPackageName converts a string to a valid Go package name.
// Package names are lowercase with no separators.
// Examples: "UserProfile" -> "userprofile", "user-profile" -> "userprofile"
func ToPackageName(s string) string {
	if s == "" {
		return ""
	}
	// Remove all non-alphanumeric characters and convert to lowercase
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

// ToVariableName converts a string to a valid Go variable name (camelCase).
// Examples: "user_profile" -> "userProfile", "UserProfile" -> "userProfile"
func ToVariableName(s string) string {
	return ToCamelCase(s)
}

// ToRepoVariableName returns the repository variable name for a domain.
// Examples: "user" -> "userRepo", "userProfile" -> "userProfileRepo"
func ToRepoVariableName(domain string) string {
	return ToCamelCase(domain) + "Repo"
}

// ToServiceVariableName returns the service variable name for a domain.
// Examples: "user" -> "userService", "userProfile" -> "userProfileService"
func ToServiceVariableName(domain string) string {
	return ToCamelCase(domain) + "Service"
}

// ToControllerVariableName returns the controller variable name for a domain.
// Examples: "user" -> "userController", "userProfile" -> "userProfileController"
func ToControllerVariableName(domain string) string {
	return ToCamelCase(domain) + "Controller"
}

// ToModelName returns the model struct name for a domain.
// Examples: "user" -> "User", "user_profile" -> "UserProfile"
func ToModelName(domain string) string {
	return ToPascalCase(domain)
}

// ToTableName returns the database table name for a domain.
// Uses plural form and snake_case.
// Examples: "user" -> "users", "user_profile" -> "user_profiles"
func ToTableName(domain string) string {
	singular := ToSnakeCase(domain)
	return inflection.Plural(singular)
}

// ToURLPath returns the URL path for a domain.
// Uses plural form and kebab-case.
// Examples: "user" -> "/users", "userProfile" -> "/user-profiles"
func ToURLPath(domain string) string {
	kebab := ToKebabCase(domain)
	plural := inflection.Plural(kebab)
	return "/" + plural
}

// Pluralize returns the plural form of a word.
// Examples: "user" -> "users", "category" -> "categories", "person" -> "people"
func Pluralize(s string) string {
	if s == "" {
		return ""
	}
	return inflection.Plural(s)
}

// Singularize returns the singular form of a word.
// Examples: "users" -> "user", "categories" -> "category", "people" -> "person"
func Singularize(s string) string {
	if s == "" {
		return ""
	}
	return inflection.Singular(s)
}

// ToJSONTag returns the JSON tag for a field name.
// Converts to snake_case by default.
// Examples: "FirstName" -> "first_name", "userID" -> "user_id"
func ToJSONTag(fieldName string) string {
	return ToSnakeCase(fieldName)
}

// ToLabel converts a field name to a human-readable label.
// Examples: "firstName" -> "First Name", "user_id" -> "User ID"
func ToLabel(s string) string {
	if s == "" {
		return ""
	}

	words := splitWords(s)
	var result []string

	for _, word := range words {
		if word == "" {
			continue
		}
		// Handle common acronyms
		upper := strings.ToUpper(word)
		if isAcronym(upper) {
			result = append(result, upper)
		} else {
			// Title case
			runes := []rune(strings.ToLower(word))
			runes[0] = unicode.ToUpper(runes[0])
			result = append(result, string(runes))
		}
	}

	return strings.Join(result, " ")
}

// splitWords splits a string into words based on common separators and case changes.
func splitWords(s string) []string {
	// First, replace common separators with spaces
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")

	// Then split on case changes (for camelCase/PascalCase)
	var words []string
	var current strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		if r == ' ' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			continue
		}

		if unicode.IsUpper(r) && i > 0 {
			// Check if this starts a new word
			prevLower := unicode.IsLower(runes[i-1])
			// Also check if next char is lowercase (handles "XMLParser" -> "XML", "Parser")
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])

			if prevLower || (i+1 < len(runes) && nextLower && current.Len() > 1) {
				if current.Len() > 0 {
					words = append(words, current.String())
					current.Reset()
				}
			}
		}
		current.WriteRune(r)
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

// isAcronym checks if a word is a common acronym.
func isAcronym(word string) bool {
	acronyms := map[string]bool{
		"ID": true, "URL": true, "URI": true, "API": true,
		"HTTP": true, "HTTPS": true, "HTML": true, "CSS": true,
		"JSON": true, "XML": true, "SQL": true, "UUID": true,
		"IP": true, "TCP": true, "UDP": true, "DNS": true,
		"CPU": true, "GPU": true, "RAM": true, "ROM": true,
		"UI": true, "UX": true, "OK": true,
	}
	return acronyms[word]
}

// ToImportAlias returns an import alias for a package path.
// Examples: "github.com/user/project/internal/services/user" -> "userService"
func ToImportAlias(pkgPath, suffix string) string {
	parts := strings.Split(pkgPath, "/")
	if len(parts) == 0 {
		return ""
	}
	base := parts[len(parts)-1]
	return ToCamelCase(base) + suffix
}

// SanitizeIdentifier removes invalid characters from an identifier.
func SanitizeIdentifier(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	for i, r := range s {
		if i == 0 {
			if unicode.IsLetter(r) || r == '_' {
				result.WriteRune(r)
			} else {
				result.WriteRune('_')
			}
		} else {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				result.WriteRune(r)
			}
		}
	}
	return result.String()
}
