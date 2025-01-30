package transformations_test

import (
    "reflect"
    "sync"
    "testing"
    
    "xml2json/transformations"
)

func TestRegisterAndGetTransformer(t *testing.T) {
    mockTransformer := func(data map[string]interface{}, param any) interface{} {
        return "hello"
    }

    transformations.Register("mock", mockTransformer)
    transformer, found := transformations.GetTransformer("mock")
    if !found {
        t.Fatalf("Expected to find transformer 'mock' after registration, but got not found")
    }
    if transformer == nil {
        t.Fatalf("Expected a non-nil transformer, but got nil")
    }
    result := transformer(nil, nil)
    if !reflect.DeepEqual(result, "hello") {
        t.Errorf("Expected transformer result = 'hello'; got %v", result)
    }
}
func TestGetTransformer_NotFound(t *testing.T) {
    transformer, found := transformations.GetTransformer("non-existent")
    if found {
        t.Errorf("Expected not to find 'non-existent' transformer, but found = %v", found)
    }
    if transformer != nil {
        t.Errorf("Expected transformer to be nil for an unregistered name; got %v", transformer)
    }
}
func TestConcurrentRegistration(t *testing.T) {
	numGoroutines := 10
    var wg sync.WaitGroup
    wg.Add(numGoroutines)

    for i := 0; i < numGoroutines; i++ {
        go func(i int) {
            defer wg.Done()

            name := "test_transformer_" + string(rune('A'+i))
            transformations.Register(name, func(data map[string]interface{}, param any) interface{} {
                return i
            })
        }(i)
    }

    wg.Wait()

    for i := 0; i < numGoroutines; i++ {
        name := "test_transformer_" + string(rune('A'+i))
        transformer, found := transformations.GetTransformer(name)
        if !found {
            t.Errorf("Expected to find transformer '%s' but not found", name)
        } else {
            val := transformer(nil, nil)
            if intVal, ok := val.(int); !ok || intVal != i {
                t.Errorf("Transformer '%s' returned unexpected value: got %v (want %d)", name, val, i)
            }
        }
    }
}
