package buildjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type SchemaTransformer struct {
	Mapping map[string]map[string]string
}

func NewSchemaTransformer() *SchemaTransformer {
	fmt.Println("NewSchemaTransformer invoked!")

	mapping := map[string]map[string]string{}

	data, err := os.ReadFile("mapping.json")

	if err != nil {
		log.Fatalf("Couldn't file mapping.json")
		return nil
	}

	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return nil
	}

	st := SchemaTransformer{
		Mapping: mapping,
	}
	return &st
}

func (st *SchemaTransformer) BuildFromJSON(jsonFilePath string) (jsonObj map[string]interface{}, err error) {
	splitPaths := strings.Split(jsonFilePath, "/")

	fileName := splitPaths[len(splitPaths)-1]

	fileSplit := strings.Split(fileName, "_")

	// path can be of the type: {provider}_{target_schema}_<extra_keywords_if_any>.json.
	// or of the type {provider}_{target_schema}.json.
	provider := fileSplit[0]
	schema := strings.Split(fileSplit[len(fileSplit)-1], ".")[0]

	schemaFilePath := fmt.Sprintf("./pkg/buildjson/schema/%s.json", schema)

	schemaData, err := os.ReadFile(schemaFilePath)
	if err != nil {
		// Wrong schema file path
		return nil, fmt.Errorf("couldn't find schema")
	}

	var schemaInput map[string]interface{}
	err = json.Unmarshal(schemaData, &schemaInput)
	if err != nil {
		return nil, err
	}

	// Get the schema data input
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		// Wrong file path
		return nil, fmt.Errorf("Wrong File path!")
	}

	var input map[string]interface{}
	err = json.Unmarshal(data, &input)
	if err != nil {
		return nil, err
	}

	// transformAndLoad(input map[stirng]interface{}) (output map[string]interface{})
	providerMap, exists := st.Mapping[provider]
	if !exists {
		return nil, fmt.Errorf("no mapping present for the provider")
	}

	jsonObj, err = st.validateAndTransform(input, providerMap)
	if err != nil {
		return nil, err
	}

	return jsonObj, nil
}

// Helper function to create nested maps from flattened keys
func (st *SchemaTransformer) nestKeys(output map[string]interface{}, fullKey string, value interface{}) {
	keys := strings.Split(fullKey, ".")
	currentMap := output

	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]

		// If the key doesn't exist, create an empty map
		if _, exists := currentMap[key]; !exists {
			currentMap[key] = make(map[string]interface{})
		}

		// Move deeper into the nested map
		currentMap = currentMap[key].(map[string]interface{})
	}

	// Assign the final key's value
	currentMap[keys[len(keys)-1]] = value
}

// Function to transform flattened input JSON based on the key mappings
func (st *SchemaTransformer) validateAndTransform(input map[string]interface{}, fieldMappings map[string]string) (map[string]interface{}, error) {
	output := make(map[string]interface{})

	for inputKey, mappedKey := range fieldMappings {
		// Split inputKey to handle nested paths
		keys := strings.Split(inputKey, ".")

		// Traverse through the input to find the value at nested path
		var currentValue interface{} = input
		for _, key := range keys {
			if currMap, ok := currentValue.(map[string]interface{}); ok {
				if value, exists := currMap[key]; exists {
					currentValue = value
				} else {
					currentValue = nil
					break
				}
			} else {
				return nil, errors.New("invalid input structure")
			}
		}

		// If a valid value exists for the inputKey, use nestKeys to add it to the output
		if currentValue != nil {
			st.nestKeys(output, mappedKey, currentValue)
		}
	}

	return output, nil
}
