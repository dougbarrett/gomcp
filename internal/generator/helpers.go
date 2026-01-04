package generator

import (
	"strings"
	"text/template"

	"github.com/dbb1dev/go-mcp/internal/utils"
)

// TemplateFuncMap returns the template function map.
func TemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		// String case conversion
		"toPascalCase": utils.ToPascalCase,
		"toCamelCase":  utils.ToCamelCase,
		"toSnakeCase":  utils.ToSnakeCase,
		"toKebabCase":  utils.ToKebabCase,
		"toLower":      strings.ToLower,
		"toUpper":      strings.ToUpper,
		"toTitle":      strings.Title, //nolint:staticcheck

		// Naming utilities
		"toPackageName":  utils.ToPackageName,
		"toVariableName": utils.ToVariableName,
		"toModelName":    utils.ToModelName,
		"toTableName":    utils.ToTableName,
		"toURLPath":      utils.ToURLPath,
		"toJSONTag":      utils.ToJSONTag,
		"toLabel":        utils.ToLabel,

		// Variable names for DI
		"toRepoVariableName":       utils.ToRepoVariableName,
		"toServiceVariableName":    utils.ToServiceVariableName,
		"toControllerVariableName": utils.ToControllerVariableName,

		// Pluralization
		"pluralize":   utils.Pluralize,
		"singularize": utils.Singularize,

		// String utilities
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"replace":    strings.ReplaceAll,
		"trim":       strings.TrimSpace,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"split":      strings.Split,
		"join":       strings.Join,

		// Comparison and logic
		"eq":  func(a, b interface{}) bool { return a == b },
		"ne":  func(a, b interface{}) bool { return a != b },
		"and": func(a, b bool) bool { return a && b },
		"or":  func(a, b bool) bool { return a || b },
		"not": func(a bool) bool { return !a },

		// Conditional helpers
		"default": func(def, val interface{}) interface{} {
			if val == nil || val == "" || val == 0 || val == false {
				return def
			}
			return val
		},
		"coalesce": func(vals ...interface{}) interface{} {
			for _, val := range vals {
				if val != nil && val != "" && val != 0 && val != false {
					return val
				}
			}
			return nil
		},

		// Type checking
		"isString": func(v interface{}) bool {
			_, ok := v.(string)
			return ok
		},
		"isBool": func(v interface{}) bool {
			_, ok := v.(bool)
			return ok
		},
		"isInt": func(v interface{}) bool {
			switch v.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				return true
			}
			return false
		},
		"isFloat": func(v interface{}) bool {
			switch v.(type) {
			case float32, float64:
				return true
			}
			return false
		},

		// Slice utilities
		"first": func(list []interface{}) interface{} {
			if len(list) > 0 {
				return list[0]
			}
			return nil
		},
		"last": func(list []interface{}) interface{} {
			if len(list) > 0 {
				return list[len(list)-1]
			}
			return nil
		},
		"len": func(v interface{}) int {
			switch val := v.(type) {
			case string:
				return len(val)
			case []interface{}:
				return len(val)
			case []string:
				return len(val)
			case map[string]interface{}:
				return len(val)
			default:
				return 0
			}
		},
		"empty": func(v interface{}) bool {
			switch val := v.(type) {
			case nil:
				return true
			case string:
				return val == ""
			case []interface{}:
				return len(val) == 0
			case []string:
				return len(val) == 0
			case map[string]interface{}:
				return len(val) == 0
			case bool:
				return !val
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				return val == 0
			case float32, float64:
				return val == 0.0
			default:
				return false
			}
		},

		// Index helpers
		"isFirst": func(i int) bool { return i == 0 },
		"isLast":  func(i, length int) bool { return i == length-1 },
		"inc":     func(i int) int { return i + 1 },
		"dec":     func(i int) int { return i - 1 },

		// Code generation helpers
		"goType": func(t string) string {
			// Map common type names to Go types
			typeMap := map[string]string{
				"text":     "string",
				"integer":  "int",
				"number":   "float64",
				"boolean":  "bool",
				"datetime": "time.Time",
				"date":     "time.Time",
				"time":     "time.Time",
				"binary":   "[]byte",
			}
			if mapped, ok := typeMap[strings.ToLower(t)]; ok {
				return mapped
			}
			return t
		},

		// templui form type mapping
		"formComponent": func(formType string) string {
			componentMap := map[string]string{
				"input":    "input.Input",
				"textarea": "textarea.Textarea",
				"select":   "select_box.SelectBox",
				"checkbox": "checkbox.Checkbox",
				"switch":   "switch_comp.Switch",
				"date":     "date_picker.DatePicker",
				"time":     "time_picker.TimePicker",
				"email":    "input.Input",
				"password": "input.Input",
				"number":   "input.Input",
				"rating":   "rating.Rating",
				"tags":     "tags_input.TagsInput",
				"slider":   "slider.Slider",
			}
			if comp, ok := componentMap[formType]; ok {
				return comp
			}
			return "input.Input"
		},

		// Input type attribute for form fields
		"inputType": func(formType string) string {
			typeMap := map[string]string{
				"email":    "email",
				"password": "password",
				"number":   "number",
				"date":     "date",
				"time":     "time",
				"datetime": "datetime-local",
			}
			if t, ok := typeMap[formType]; ok {
				return t
			}
			return "text"
		},

		// Faker function for seeding
		"fakerFunc": func(goType string) string {
			fakerMap := map[string]string{
				"string":     "gofakeit.Word()",
				"int":        "gofakeit.Number(1, 1000)",
				"int64":      "int64(gofakeit.Number(1, 1000))",
				"uint":       "uint(gofakeit.Number(1, 1000))",
				"float64":    "gofakeit.Float64Range(0, 1000)",
				"bool":       "gofakeit.Bool()",
				"time.Time":  "gofakeit.Date()",
				"*time.Time": "func() *time.Time { t := gofakeit.Date(); return &t }()",
			}
			if fn, ok := fakerMap[goType]; ok {
				return fn
			}
			return "nil"
		},

		// GORM tag helpers
		"gormTag": func(fieldType string, required bool, tags string) string {
			var parts []string

			// Add size for strings
			if fieldType == "string" && !strings.Contains(tags, "size:") && !strings.Contains(tags, "type:text") {
				parts = append(parts, "size:255")
			}

			// Add not null for required
			if required && !strings.Contains(tags, "not null") {
				parts = append(parts, "not null")
			}

			// Add custom tags
			if tags != "" {
				parts = append(parts, tags)
			}

			if len(parts) == 0 {
				return ""
			}
			return strings.Join(parts, ";")
		},

		// Comment helpers
		"comment": func(s string) string {
			if s == "" {
				return ""
			}
			return "// " + s
		},
		"blockComment": func(s string) string {
			if s == "" {
				return ""
			}
			return "/* " + s + " */"
		},

		// Indent helper
		"indent": func(spaces int, s string) string {
			pad := strings.Repeat(" ", spaces)
			lines := strings.Split(s, "\n")
			for i, line := range lines {
				if line != "" {
					lines[i] = pad + line
				}
			}
			return strings.Join(lines, "\n")
		},

		// Quote helpers
		"quote":       func(s string) string { return `"` + s + `"` },
		"singleQuote": func(s string) string { return `'` + s + `'` },
		"backtick":    func(s string) string { return "`" + s + "`" },

		// Dict helper for passing multiple values to templates
		"dict": func(values ...interface{}) map[string]interface{} {
			if len(values)%2 != 0 {
				return nil
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					continue
				}
				dict[key] = values[i+1]
			}
			return dict
		},

		// List helper
		"list": func(values ...interface{}) []interface{} {
			return values
		},
	}
}
