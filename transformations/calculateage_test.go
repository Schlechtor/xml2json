package transformations_test

import (
	"reflect"
	"testing"

	"xml2json/transformations"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		param   interface{}
		wantNil bool
		wantAge interface{}
	}{
		{
			name:    "param not a string",
			data:    map[string]interface{}{"DateOfBirth": "2000-01-01"},
			param:   1234,
			wantNil: true,
		},
		{
			name:    "empty param string",
			data:    map[string]interface{}{"DateOfBirth": "2000-01-01"},
			param:   "",
			wantNil: true,
		},
		{
			name:    "field doesn't exist in data",
			data:    map[string]interface{}{"SomethingElse": "2020-02-20"},
			param:   "DateOfBirth",
			wantNil: true,
		},
		{
			name:    "invalid date format",
			data:    map[string]interface{}{"DateOfBirth": "85-07-15"},
			param:   "DateOfBirth",
			wantNil: true,
		},
		{
			name:    "valid date, check approximate age",
			data:    map[string]interface{}{"DateOfBirth": "1985-07-15"},
			param:   "DateOfBirth",
			wantNil: false,
			wantAge: 39,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := transformations.GetTransformer("calculate")
			if got == nil {
				t.Fatal("Expected 'calculate' transformer to be registered but got nil")
			}
			result := got(tt.data, tt.param)

			if tt.wantNil {
				if result != nil {
					t.Errorf("Expected result to be nil, but got: %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("Expected non-nil result, got nil")
				} else {
					if !reflect.DeepEqual(result, tt.wantAge) {
						t.Errorf("Expected age %v, got %v", tt.wantAge, result)
					}
				}
			}
		})
	}
}
