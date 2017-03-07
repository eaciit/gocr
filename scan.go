package gocr

import (
	"encoding/gob"
	"os"
)

// Read the model from a file and return Model
func ReadModel(path string) Model {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	// Close the file later
	defer file.Close()

	// Initialize decoder and empty Model
	decoder := gob.NewDecoder(file)
	model := Model{}

	// Decode the file
	decoder.Decode(&model)

	return model
}
