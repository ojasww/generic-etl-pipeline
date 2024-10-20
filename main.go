package main

import (
	"encoding/json"
	"etl/pkg/buildjson"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	builder := buildjson.NewSchemaTransformer()

	folderPath := "./sample_data/input_files"
	outputFolderPath := "./sample_data/output_files"

	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())
		fmt.Println(filePath)

		jsonObj, err := builder.BuildFromJSON(filePath)
		if err != nil {
			log.Fatalf("Error building the JSON", err)
		}
		// fmt.Println("jsonObj", jsonObj)

		// Convert the jsonObj to JSON string for writing
		jsonData, err := json.MarshalIndent(jsonObj, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling JSON object: %v", err)
		}

		// Determine output file path
		outputFileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())) + "_output.json"
		outputFilePath := filepath.Join(outputFolderPath, outputFileName)

		// Write the JSON object to the output file
		err = os.WriteFile(outputFilePath, jsonData, 0644)
		if err != nil {
			log.Fatalf("Error writing JSON file: %v", err)
		}

	}

	if err != nil {
		log.Fatalf("Error walking the path: %v", err)
	}
}
