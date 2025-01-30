package transformations

import (
	"fmt"
	"strings"
)

func init() {
	Register("combine", combineFields)
}

func combineFields(data map[string]interface{}, param any) interface{} {
	fields, ok := param.([]interface{})
	if !ok {
		return nil
	}
	var combinedValue strings.Builder
	for _, field := range fields {
		fieldStr, ok := field.(string)
		if ok && fieldStr != "" {
			if fieldValue, exists := data[fieldStr]; exists {
				combinedValue.WriteString(" ")
				combinedValue.WriteString(fmt.Sprintf("%v", fieldValue))
			}
		}
	}
	return strings.TrimSpace(combinedValue.String())
}
