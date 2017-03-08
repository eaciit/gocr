package gocr

import (
	"encoding/csv"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver"
	"github.com/ugorji/go/codec"
)

type ModelImage struct {
	Label string
	Data  ImageMatrix
}

type Model struct {
	Name        string
	ModelImages []ModelImage
}

// Download the dataset from here and save it in chars74k_dataset/EnglishFnt
// Also copy index.csv to extracted folder in targetPath
func Prepare(targetPath string) error {
	// Url for dataset of computer fonts image
	url := "http://www.ee.surrey.ac.uk/CVSSP/demos/chars74k/EnglishFnt.tgz"
	filePath := targetPath + "english_dataset.tgz"

	fmt.Println("Downloading english font dataset...")

	// Downloading dataset file
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create the output file
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Save to file
	n, err := io.Copy(outFile, response.Body)
	if err != nil {
		return err
	}

	// Done Downloading
	fmt.Println(n, " bytes of english dataset successfully downloaded.")

	// Extract the file
	err = archiver.TarGz.Open(filePath, targetPath)
	if err != nil {
		return err
	}

	extractedPath := targetPath + "English/Fnt/"
	fmt.Println("Dataset extracted to ", extractedPath)

	// Get packagePath
	_, b, _, _ := runtime.Caller(0)
	packagePath := filepath.Dir(b)

	// Open provided index.csv
	srcFile, err := os.Open(packagePath + "/train_data/chars74k_dataset/index.csv")
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination index.csv
	destFile, err := os.Create(extractedPath + "index.csv") // creates if file doesn't exist
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy index.csv of the chars74k dataset to target folder
	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func ReadCSV(path string) ([][]string, error) {
	// Read the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Close the file later
	defer file.Close()

	reader := csv.NewReader(file)
	var datas [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		datas = append(datas, record)
	}

	return datas, nil
}

// Training method
// The train folder should include index.csv and images that inside the index.csv
func Train(sampleFolderPath string, modelPath string) error {

	indexPath := sampleFolderPath + "/index.csv"
	indexData, err := ReadCSV(indexPath)
	if err != nil {
		return err
	}

	// Initialize model
	model := Model{}

	// Read and binarize each image to array 0 and 1
	for _, elm := range indexData {
		image, readErr := ReadImage(sampleFolderPath + elm[0])
		if err != nil {
			return readErr
		}

		binaryArray := ImageToBinaryArray(image)

		model.ModelImages = append(model.ModelImages, ModelImage{
			Label: elm[1],
			Data:  binaryArray,
		})
	}

	// Create the model file
	modelFile, err := os.Create(modelPath + "model.gob")
	if err != nil {
		return err
	}

	// Close the file later
	defer modelFile.Close()

	// Create encoder
	encoder := codec.NewEncoder(modelFile, new(codec.CborHandle))

	// Write the file
	if err := encoder.Encode(model); err != nil {
		return err
	}

	return nil
}
