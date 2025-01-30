package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"xml2json/transformations"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TopLevelField string                 `yaml:"topLevelField"`
	Fields        map[string]interface{} `yaml:"fields"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func processXML(xmlData []byte, config *Config) ([]map[string]interface{}, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(xmlData)))
	var jsons []map[string]interface{}
	var current map[string]interface{}
	var field string

	topLevelField := config.TopLevelField

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch ele := tok.(type) {
		case xml.StartElement:
			if ele.Name.Local == topLevelField {
				current = make(map[string]interface{})
				for _, attr := range ele.Attr {
					current[attr.Name.Local] = attr.Value
				}
			} else if current != nil {
				field = ele.Name.Local
				value := decodeXMLElement(decoder)
				current[field] = value
			}

		case xml.EndElement:
			if ele.Name.Local == topLevelField && current != nil {
				transformed := Transform(current, config)
				jsons = append(jsons, transformed)
			}
		}
	}

	if len(jsons) == 0 {
		return nil, fmt.Errorf("no data found in the XML")
	}

	return jsons, nil
}

func decodeXMLElement(decoder *xml.Decoder) interface{} {
	result := make(map[string]interface{})
	var charData strings.Builder

	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		switch t := tok.(type) {
		case xml.StartElement:
			result[t.Name.Local] = decodeXMLElement(decoder)
		case xml.CharData:
			charData.WriteString(strings.TrimSpace(string(t)))
		case xml.EndElement:
			if len(result) > 0 {
				return result
			}
			return charData.String()
		}
	}

	return result
}

func Transform(data map[string]interface{}, config *Config) map[string]interface{} {
	transformed := make(map[string]interface{})

	for key, transformation := range config.Fields {
		if t, ok := transformation.(map[interface{}]interface{}); ok {
			strMap := make(map[string]interface{})
			for k, v := range t {
				if strKey, ok := k.(string); ok {
					strMap[strKey] = v
				}
			}

			var handled bool
			for operation, params := range strMap {
				if handler, found := transformations.GetTransformer(operation); found {
					val := handler(data, params)
					if val != nil {
						transformed[key] = val
					}
					handled = true
					break
				}
			}

			if !handled {
				if fieldName, ok := strMap["field"].(string); ok {
					fieldValue := extractFieldData(data, fieldName)
					if fieldValue == nil {
						continue
					}

					if nestedConf, ok := strMap["subfields"].(map[string]interface{}); ok {
						nestedVal := extractNestedFields(fieldValue, nestedConf)
						if len(nestedVal) > 0 {
							transformed[key] = nestedVal
						}
					} else {
						transformed[key] = fieldValue
					}
				}
			}
		} else if field, ok := transformation.(string); ok {
			if value, exists := data[field]; exists {
				transformed[key] = value
			}
		}
	}

	return transformed
}

func extractFieldData(data map[string]interface{}, fieldName string) interface{} {
	if fieldName == "" {
		return nil
	}
	if value, exists := data[fieldName]; exists {
		return value
	}
	return nil
}

func extractNestedFields(fieldData interface{}, subfields map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	if nestedData, ok := fieldData.(map[string]interface{}); ok {
		for outKey, subfieldRaw := range subfields {
			if subfieldStr, ok := subfieldRaw.(string); ok {
				if val, exists := nestedData[subfieldStr]; exists {
					result[outKey] = val
				}
			}
		}
	}

	return result
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <config.yaml> <input.xml>")
		return
	}

	configFile := os.Args[1]
	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	xmlFile := os.Args[2]
	xmlData, err := os.ReadFile(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML file:", err)
		return
	}

	jsons, err := processXML(xmlData, config)
	if err != nil {
		fmt.Println("Error processing XML:", err)
		return
	}

	output := map[string]interface{}{
		config.TopLevelField: jsons,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
