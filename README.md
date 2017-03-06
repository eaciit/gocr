# gocr
gocr is a go based OCR module

# How to use
## Train Data
To use OCR a trained model need to be build first using Train method. Train read sample of many files from a path. These path will contain many image file with only 2 colors (black over white) and an index file (index.csv) which contains key value to identify file name and characted represent by the file

ie:
a.gif,a
b.gif,b
A1.gif,A
B1.gif,B

```go
pathOfSample := "/usr/eaciit/go/src/github.com/eaciit/gocr/sample"
pathOfTrainedModel := "/usr/eaciit/go/src/github.com/eaciit/gocr/trainedmodel"
ocr := new Ocr()
if trainResult, err := ocr.Train(pathOfSample, pathOfTrainedModel); err!=nil {
  fmt.Println("Error :" + err.Error())
} else {
  fmt.Prinln("Trained ",trainResult)
}
```

## Scan based on Trained Model
```go
pathOfTrainedModel := "/usr/eaciit/go/src/github.com/eaciit/gocr/trainedmodel"
ocr := new Ocr()
ocr.TrainedModel = pathOfTrainedModel
if scanResult, err := ocr.Scan("/usr/eaciit/doc1.pdf"); err!=nil {
  fmt.Println("Error :" + err.Error())
} else {
  fmt.Prinln("Scanned: ",scanResult)
}
```

