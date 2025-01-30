package transformations_test

import (
    "reflect"
    "testing"

    "xml2json/transformations"
)

func TestCombineFields(t *testing.T) {
    transformer, found := transformations.GetTransformer("combine")
    if !found {
        t.Fatal("Expected 'combine' transformer to be registered, but not found")
    }
    if transformer == nil {
        t.Fatal("Expected a non-nil transformer for 'combine', but got nil")
    }

    tests := []struct {
        name   string
        data   map[string]interface{}
        param  interface{}
        want   interface{}
    }{
        {
            name:  "Param not a slice => returns nil",
            data:  map[string]interface{}{"FirstName": "John", "LastName": "Doe"},
            param: "invalid param",
            want:  nil,
        },
        {
            name:  "Param is an empty slice => returns empty string",
            data:  map[string]interface{}{"FirstName": "John", "LastName": "Doe"},
            param: []interface{}{},
            want:  "",
        },
        {
            name:  "Fields exist => combine them with space",
            data:  map[string]interface{}{"FirstName": "John", "LastName": "Doe"},
            param: []interface{}{"FirstName", "LastName"},
            want:  "John Doe",
        },
        {
            name:  "Some fields missing => skip missing ones",
            data:  map[string]interface{}{"FirstName": "John"},
            param: []interface{}{"FirstName", "LastName"},
            // "LastName" not in data, so only "John"
            want:  "John",
        },
        {
            name:  "Non-string field in param => skip it",
            data:  map[string]interface{}{"FirstName": "John", "LastName": "Doe"},
            param: []interface{}{"FirstName", 9999, "LastName"},
            want:  "John Doe", 
        },
        {
            name:  "Empty strings in param => skip them",
            data:  map[string]interface{}{"FirstName": "John", "LastName": "Doe"},
            param: []interface{}{"FirstName", "", "LastName"},
            want:  "John Doe",
        },
        {
            name: "Fields with different data types => still appended as string",
            data: map[string]interface{}{
                "Age":     30,
                "City":    "Berlin",
                "IsValid": true,
            },
            param: []interface{}{"Age", "City", "IsValid"},
            want: "30 Berlin true",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := transformer(tt.data, tt.param)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("combineFields(%v, %v) = %v; want %v", tt.data, tt.param, got, tt.want)
            }
        })
    }
}
