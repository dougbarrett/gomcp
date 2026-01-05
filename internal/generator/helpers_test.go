package generator

import (
	"testing"
)

// TestTemplateFuncMap tests that all expected functions are present.
func TestTemplateFuncMap(t *testing.T) {
	funcMap := TemplateFuncMap()

	expectedFuncs := []string{
		// String case conversion
		"toPascalCase", "toCamelCase", "toSnakeCase", "toKebabCase",
		"toLower", "toUpper", "toTitle",
		// Naming utilities
		"toPackageName", "toVariableName", "toModelName", "toTableName",
		"toURLPath", "toJSONTag", "toLabel",
		// DI variable names
		"toRepoVariableName", "toServiceVariableName", "toControllerVariableName",
		// Pluralization
		"pluralize", "singularize",
		// String utilities
		"contains", "hasPrefix", "hasSuffix", "replace", "trim",
		"trimPrefix", "trimSuffix", "split", "join",
		// Comparison and logic
		"eq", "ne", "and", "or", "not",
		// Conditional helpers
		"default", "coalesce",
		// Type checking
		"isString", "isBool", "isInt", "isFloat",
		// Slice utilities
		"first", "last", "len", "empty",
		// Index helpers
		"isFirst", "isLast", "inc", "dec",
		// Code generation helpers
		"goType", "formComponent", "inputType", "fakerFunc", "gormTag",
		// Comment helpers
		"comment", "blockComment",
		// Indent helper
		"indent",
		// Quote helpers
		"quote", "singleQuote", "backtick",
		// Dict and list
		"dict", "list",
	}

	for _, fn := range expectedFuncs {
		if _, ok := funcMap[fn]; !ok {
			t.Errorf("TemplateFuncMap missing function: %s", fn)
		}
	}
}

// TestHelperFunctions_CaseConversion tests case conversion functions.
func TestHelperFunctions_CaseConversion(t *testing.T) {
	tests := []struct {
		content string
		data    map[string]string
		want    string
	}{
		{"[[.X | toPascalCase]]", map[string]string{"X": "user_profile"}, "UserProfile"},
		{"[[.X | toCamelCase]]", map[string]string{"X": "user_profile"}, "userProfile"},
		{"[[.X | toSnakeCase]]", map[string]string{"X": "UserProfile"}, "user_profile"},
		{"[[.X | toKebabCase]]", map[string]string{"X": "UserProfile"}, "user-profile"},
		{"[[.X | toLower]]", map[string]string{"X": "HELLO"}, "hello"},
		{"[[.X | toUpper]]", map[string]string{"X": "hello"}, "HELLO"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_Naming tests naming utility functions.
func TestHelperFunctions_Naming(t *testing.T) {
	tests := []struct {
		content string
		data    map[string]string
		want    string
	}{
		{"[[.X | toPackageName]]", map[string]string{"X": "UserProfile"}, "userprofile"},
		{"[[.X | toVariableName]]", map[string]string{"X": "user_profile"}, "userProfile"},
		{"[[.X | toModelName]]", map[string]string{"X": "user_profile"}, "UserProfile"},
		{"[[.X | toTableName]]", map[string]string{"X": "UserProfile"}, "user_profiles"},
		{"[[.X | toURLPath]]", map[string]string{"X": "UserProfile"}, "/user-profiles"},
		{"[[.X | toJSONTag]]", map[string]string{"X": "UserName"}, "user_name"},
		{"[[.X | toLabel]]", map[string]string{"X": "user_name"}, "User Name"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_DIVariables tests DI variable name functions.
func TestHelperFunctions_DIVariables(t *testing.T) {
	tests := []struct {
		content string
		data    map[string]string
		want    string
	}{
		{"[[.X | toRepoVariableName]]", map[string]string{"X": "user"}, "userRepo"},
		{"[[.X | toServiceVariableName]]", map[string]string{"X": "user"}, "userService"},
		{"[[.X | toControllerVariableName]]", map[string]string{"X": "user"}, "userController"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_Pluralization tests pluralize/singularize functions.
func TestHelperFunctions_Pluralization(t *testing.T) {
	tests := []struct {
		content string
		data    map[string]string
		want    string
	}{
		{"[[.X | pluralize]]", map[string]string{"X": "user"}, "users"},
		{"[[.X | pluralize]]", map[string]string{"X": "category"}, "categories"},
		{"[[.X | singularize]]", map[string]string{"X": "users"}, "user"},
		{"[[.X | singularize]]", map[string]string{"X": "categories"}, "category"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_StringUtils tests string utility functions.
func TestHelperFunctions_StringUtils(t *testing.T) {
	tests := []struct {
		content string
		data    map[string]string
		want    string
	}{
		{`[[if contains .X "llo"]]yes[[end]]`, map[string]string{"X": "hello"}, "yes"},
		{`[[if hasPrefix .X "he"]]yes[[end]]`, map[string]string{"X": "hello"}, "yes"},
		{`[[if hasSuffix .X "lo"]]yes[[end]]`, map[string]string{"X": "hello"}, "yes"},
		{`[[replace .X "l" "L"]]`, map[string]string{"X": "hello"}, "heLLo"},
		{`[[trim .X]]`, map[string]string{"X": "  hello  "}, "hello"},
		{`[[trimPrefix .X "he"]]`, map[string]string{"X": "hello"}, "llo"},
		{`[[trimSuffix .X "lo"]]`, map[string]string{"X": "hello"}, "hel"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_Comparison tests comparison functions.
func TestHelperFunctions_Comparison(t *testing.T) {
	tests := []struct {
		content string
		data    interface{}
		want    string
	}{
		{`[[if eq .A .B]]same[[end]]`, map[string]string{"A": "x", "B": "x"}, "same"},
		{`[[if ne .A .B]]diff[[end]]`, map[string]string{"A": "x", "B": "y"}, "diff"},
		{`[[if and .A .B]]both[[end]]`, map[string]bool{"A": true, "B": true}, "both"},
		{`[[if or .A .B]]either[[end]]`, map[string]bool{"A": false, "B": true}, "either"},
		{`[[if not .A]]negated[[end]]`, map[string]bool{"A": false}, "negated"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_Default tests default and coalesce functions.
func TestHelperFunctions_Default(t *testing.T) {
	tests := []struct {
		content string
		data    interface{}
		want    string
	}{
		{`[[default "fallback" .X]]`, map[string]string{"X": ""}, "fallback"},
		{`[[default "fallback" .X]]`, map[string]string{"X": "value"}, "value"},
		{`[[coalesce .A .B .C]]`, map[string]string{"A": "", "B": "", "C": "third"}, "third"},
		{`[[coalesce .A .B .C]]`, map[string]string{"A": "", "B": "second", "C": "third"}, "second"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_TypeChecking tests type checking functions.
func TestHelperFunctions_TypeChecking(t *testing.T) {
	tests := []struct {
		content string
		data    interface{}
		want    string
	}{
		{`[[if isString .X]]str[[end]]`, map[string]interface{}{"X": "hello"}, "str"},
		{`[[if isBool .X]]bool[[end]]`, map[string]interface{}{"X": true}, "bool"},
		{`[[if isInt .X]]int[[end]]`, map[string]interface{}{"X": 42}, "int"},
		{`[[if isFloat .X]]float[[end]]`, map[string]interface{}{"X": 3.14}, "float"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_IndexHelpers tests index helper functions.
func TestHelperFunctions_IndexHelpers(t *testing.T) {
	tests := []struct {
		content string
		data    interface{}
		want    string
	}{
		{`[[if isFirst .I]]first[[end]]`, map[string]int{"I": 0}, "first"},
		{`[[if isLast .I .L]]last[[end]]`, map[string]int{"I": 4, "L": 5}, "last"},
		{`[[inc .I]]`, map[string]int{"I": 5}, "6"},
		{`[[dec .I]]`, map[string]int{"I": 5}, "4"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_GoType tests goType mapping.
func TestHelperFunctions_GoType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"text", "string"},
		{"integer", "int"},
		{"number", "float64"},
		{"boolean", "bool"},
		{"datetime", "time.Time"},
		{"date", "time.Time"},
		{"time", "time.Time"},
		{"binary", "[]byte"},
		{"string", "string"},         // already Go type
		{"CustomType", "CustomType"}, // unknown type passthrough
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			content := `[[goType .Type]]`
			data := map[string]string{"Type": tt.input}
			result, err := ExecuteTemplateString("test", content, data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("goType(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_FormComponent tests formComponent mapping.
func TestHelperFunctions_FormComponent(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"input", "components.Input"},
		{"textarea", "components.Textarea"},
		{"select", "components.Select"},
		{"checkbox", "components.Checkbox"},
		{"date", "components.Input"},
		{"time", "components.Input"},
		{"email", "components.Input"},
		{"password", "components.Input"},
		{"number", "components.Input"},
		{"unknown", "components.Input"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			content := `[[formComponent .Type]]`
			data := map[string]string{"Type": tt.input}
			result, err := ExecuteTemplateString("test", content, data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("formComponent(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_InputType tests inputType mapping.
func TestHelperFunctions_InputType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"email", "email"},
		{"password", "password"},
		{"number", "number"},
		{"date", "date"},
		{"time", "time"},
		{"datetime", "datetime-local"},
		{"text", "text"},    // default
		{"unknown", "text"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			content := `[[inputType .Type]]`
			data := map[string]string{"Type": tt.input}
			result, err := ExecuteTemplateString("test", content, data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("inputType(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_FakerFunc tests fakerFunc mapping.
func TestHelperFunctions_FakerFunc(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"string", "gofakeit.Word()"},
		{"int", "gofakeit.Number(1, 1000)"},
		{"int64", "int64(gofakeit.Number(1, 1000))"},
		{"uint", "uint(gofakeit.Number(1, 1000))"},
		{"float64", "gofakeit.Float64Range(0, 1000)"},
		{"bool", "gofakeit.Bool()"},
		{"time.Time", "gofakeit.Date()"},
		{"*time.Time", "func() *time.Time { t := gofakeit.Date(); return &t }()"},
		{"unknown", "nil"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			content := `[[fakerFunc .Type]]`
			data := map[string]string{"Type": tt.input}
			result, err := ExecuteTemplateString("test", content, data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("fakerFunc(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_GormTag tests gormTag function.
func TestHelperFunctions_GormTag(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		data      interface{}
		wantEmpty bool
		contains  string
	}{
		{
			name:     "string type adds size",
			content:  `[[gormTag "string" false ""]]`,
			data:     nil,
			contains: "size:255",
		},
		{
			name:     "required adds not null",
			content:  `[[gormTag "string" true ""]]`,
			data:     nil,
			contains: "not null",
		},
		{
			name:     "custom tags passed through",
			content:  `[[gormTag "string" false "index"]]`,
			data:     nil,
			contains: "index",
		},
		{
			name:      "non-string type no size",
			content:   `[[gormTag "int" false ""]]`,
			data:      nil,
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if tt.wantEmpty {
				if result != "" {
					t.Errorf("Expected empty string, got %q", result)
				}
			} else if tt.contains != "" {
				if result == "" || !containsString(result, tt.contains) {
					t.Errorf("Result %q should contain %q", result, tt.contains)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestHelperFunctions_Comment tests comment helper functions.
func TestHelperFunctions_Comment(t *testing.T) {
	tests := []struct {
		content string
		data    map[string]string
		want    string
	}{
		{`[[comment .X]]`, map[string]string{"X": "hello"}, "// hello"},
		{`[[comment .X]]`, map[string]string{"X": ""}, ""},
		{`[[blockComment .X]]`, map[string]string{"X": "hello"}, "/* hello */"},
		{`[[blockComment .X]]`, map[string]string{"X": ""}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_Indent tests indent helper function.
func TestHelperFunctions_Indent(t *testing.T) {
	content := `[[indent 4 .X]]`
	data := map[string]string{"X": "line1\nline2\nline3"}
	want := "    line1\n    line2\n    line3"

	result, err := ExecuteTemplateString("test", content, data)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if result != want {
		t.Errorf("Got %q, want %q", result, want)
	}
}

// TestHelperFunctions_Quote tests quote helper functions.
func TestHelperFunctions_Quote(t *testing.T) {
	tests := []struct {
		content string
		data    map[string]string
		want    string
	}{
		{`[[quote .X]]`, map[string]string{"X": "hello"}, `"hello"`},
		{`[[singleQuote .X]]`, map[string]string{"X": "hello"}, `'hello'`},
		{"[[backtick .X]]", map[string]string{"X": "hello"}, "`hello`"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_Dict tests dict helper function.
func TestHelperFunctions_Dict(t *testing.T) {
	content := `[[with dict "name" "John" "age" 30]][[ .name ]] is [[ .age ]][[end]]`
	result, err := ExecuteTemplateString("test", content, nil)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	want := "John is 30"
	if result != want {
		t.Errorf("Got %q, want %q", result, want)
	}
}

// TestHelperFunctions_List tests list helper function.
func TestHelperFunctions_List(t *testing.T) {
	content := `[[range list "a" "b" "c"]][[.]][[end]]`
	result, err := ExecuteTemplateString("test", content, nil)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	want := "abc"
	if result != want {
		t.Errorf("Got %q, want %q", result, want)
	}
}

// TestHelperFunctions_Empty tests empty helper function.
func TestHelperFunctions_Empty(t *testing.T) {
	tests := []struct {
		name    string
		content string
		data    interface{}
		want    string
	}{
		{"nil", `[[if empty .X]]empty[[end]]`, map[string]interface{}{"X": nil}, "empty"},
		{"empty string", `[[if empty .X]]empty[[end]]`, map[string]interface{}{"X": ""}, "empty"},
		{"empty slice", `[[if empty .X]]empty[[end]]`, map[string]interface{}{"X": []string{}}, "empty"},
		{"non-empty", `[[if empty .X]]empty[[else]]full[[end]]`, map[string]interface{}{"X": "value"}, "full"},
		{"zero int", `[[if empty .X]]empty[[end]]`, map[string]interface{}{"X": 0}, "empty"},
		{"false bool", `[[if empty .X]]empty[[end]]`, map[string]interface{}{"X": false}, "empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}

// TestHelperFunctions_Len tests len helper function.
func TestHelperFunctions_Len(t *testing.T) {
	tests := []struct {
		name    string
		content string
		data    interface{}
		want    string
	}{
		{"string", `[[len .X]]`, map[string]interface{}{"X": "hello"}, "5"},
		{"slice", `[[len .X]]`, map[string]interface{}{"X": []string{"a", "b", "c"}}, "3"},
		{"empty", `[[len .X]]`, map[string]interface{}{"X": ""}, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExecuteTemplateString("test", tt.content, tt.data)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if result != tt.want {
				t.Errorf("Got %q, want %q", result, tt.want)
			}
		})
	}
}
