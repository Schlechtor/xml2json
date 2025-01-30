package transformations

import (
	"fmt"
	"time"
)

func init() {
	Register("calculate", calculateAge)
}

func calculateAge(data map[string]interface{}, param any) interface{} {
	fieldName, ok := param.(string)
	if !ok || fieldName == "" {
		return nil
	}
	if dob, exists := data[fieldName]; exists {
		layout := "2006-01-02"
		birthDate, err := time.Parse(layout, fmt.Sprintf("%v", dob))
		if err != nil {
			return nil
		}

		currentDate := time.Now()
		age := currentDate.Year() - birthDate.Year()

		if currentDate.YearDay() < birthDate.YearDay() {
			age--
		}

		return age
	}
	return nil
}
