package itinerary

import (
	"fmt"
	"html/template"
	"strings"
)

// TemplateFuncs returns a map of template functions for Gin
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// Math operations
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"divide": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"multiply": func(a, b float64) float64 {
			return a * b
		},

		// Comparison operations
		"gt": func(a, b interface{}) bool {
			aInt := toInt(a)
			bInt := toInt(b)
			return aInt > bInt
		},
		"lt": func(a, b interface{}) bool {
			aInt := toInt(a)
			bInt := toInt(b)
			return aInt < bInt
		},
		"gte": func(a, b interface{}) bool {
			aInt := toInt(a)
			bInt := toInt(b)
			return aInt >= bInt
		},
		"lte": func(a, b interface{}) bool {
			aInt := toInt(a)
			bInt := toInt(b)
			return aInt <= bInt
		},
		"eq": func(a, b interface{}) bool {
			aInt := toInt(a)
			bInt := toInt(b)
			return aInt == bInt
		},
		"ne": func(a, b interface{}) bool {
			aInt := toInt(a)
			bInt := toInt(b)
			return aInt != bInt
		},

		// String operations
		"toUpper": func(s string) string {
			return strings.ToUpper(s)
		},
		"toLower": func(s string) string {
			return strings.ToLower(s)
		},

		// Type icons
		"typeIcon": func(typeStr string) string {
			icons := map[string]string{
				"stay":      "🏨",
				"food":      "🍽️",
				"activity":  "🎯",
				"transport": "🚗",
				"other":     "📌",
			}
			if icon, ok := icons[typeStr]; ok {
				return icon
			}
			return "📌"
		},

		// Currency formatting
		"formatPrice": func(price float64) string {
			// Format with 0 decimal places
			return "₹" + fmt.Sprintf("%.0f", price)
		},

		// Truncate text
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},

		// Format float with decimals
		"formatFloat": func(value float64, decimals int) string {
			if decimals == 0 {
				return fmt.Sprintf("%.0f", value)
			}
			formatStr := fmt.Sprintf("%%.%df", decimals)
			return fmt.Sprintf(formatStr, value)
		},
	}
}

// toInt safely converts interface{} to int
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		var i int
		fmt.Sscanf(val, "%d", &i)
		return i
	default:
		return 0
	}
}
