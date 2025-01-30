package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
    tempFile, err := ioutil.TempFile("", "config-*.yaml")
    if err != nil {
        t.Fatal("Failed to create temp file:", err)
    }
    defer os.Remove(tempFile.Name())

    yamlContent := `
topLevelField: "Patient"
fields:
  ID: "ID"
  firstName: "FirstName"
`
    if _, err := tempFile.Write([]byte(yamlContent)); err != nil {
        t.Fatal("Failed to write to temp file:", err)
    }
    tempFile.Close()

    config, err := LoadConfig(tempFile.Name())
    if err != nil {
        t.Fatalf("LoadConfig returned error for valid YAML: %v", err)
    }

    if config.TopLevelField != "Patient" {
        t.Errorf("Expected topLevelField = 'Patient', got %q", config.TopLevelField)
    }
    if len(config.Fields) != 2 {
        t.Errorf("Expected 2 fields in config, got %d", len(config.Fields))
    }

    _, err = LoadConfig("nonexistent-file.yaml")
    if err == nil {
        t.Error("Expected an error when loading non-existent file, got nil")
    }
}

func TestProcessXML(t *testing.T) {
    config := &Config{
        TopLevelField: "Patient",
        Fields: map[string]interface{}{
            "id":       "ID",
            "fname":    "FirstName",
            "lname":    "LastName",
        },
    }

    xmlString := `
<Patients>
  <Patient ID="123">
    <FirstName>John</FirstName>
    <LastName>Doe</LastName>
  </Patient>
  <Patient ID="456">
    <FirstName>Jane</FirstName>
    <LastName>Smith</LastName>
  </Patient>
</Patients>
`
    results, err := processXML([]byte(xmlString), config)
    if err != nil {
        t.Fatalf("processXML returned error: %v", err)
    }
    if len(results) != 2 {
        t.Fatalf("Expected 2 Patients, got %d", len(results))
    }
    first := results[0]
    if first["id"] != "123" {
        t.Errorf("Expected first patient ID=123, got %v", first["id"])
    }
    if first["fname"] != "John" {
        t.Errorf("Expected first patient FirstName=John, got %v", first["fname"])
    }
    second := results[1]
    if second["id"] != "456" {
        t.Errorf("Expected second patient ID=456, got %v", second["id"])
    }
    if second["lname"] != "Smith" {
        t.Errorf("Expected second patient LastName=Smith, got %v", second["lname"])
    }
}

func TestTransform(t *testing.T) {
    config := &Config{
        Fields: map[string]interface{}{
            "id": "ID",
        },
    }
    data := map[string]interface{}{
        "ID":   "ABC123",
        "Extra": "IgnoreMe",
    }

    got := Transform(data, config)
    want := map[string]interface{}{
        "id": "ABC123",
    }

    if !reflect.DeepEqual(got, want) {
        t.Errorf("Transform() got %v, want %v", got, want)
    }
}
