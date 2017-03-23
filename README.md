# gocr
gocr is a go based OCR module

## Available Scanner
* Linear Mean Scanner
* Circular Scanner (Currently the best)

## Available predictor
* NNPredictor (kNN with k = 1)
* CNNPredictor (Need to install tensorflow first) (Doesn't upport custom train)

## Future Improvement
* Sauvola Segmentation
* Better character detection
* Better symbol detection
* Spelling Correction

# How to use

To use OCR we have some generated Model that we train using tensorflow. We use Convolutional Neural Network as the model architecture and produce good detection.

You can load and use the model to `CNN Predictor` like this example below:

```go
d, _ := os.Getwd()

image, _ := gocr.ReadImage(d + "/imagetext_3.png")
s := gocr.NewCNNPredictorFromDir(modelPath + "tensor_4/")

// Define the image size
s.InputHeight, s.InputWidth = 64, 64

strings := gocr.ScanToStrings(s, image)
for _, s := range strings {
  fmt.Println(s)
}
```

However you can also use your own train data. Currently the predictor that support custom training only `NNPredictor`. Training takes `csv` file that have file image path and string representation.

ie:
```
a.gif,a
b.gif,b
A1.gif,A
B1.gif,B
```

Train the sample data to model and scan image
```go
d, _ := os.Getwd()

// Load the sample data and save it in model.cbor file
err := gocr.TrainAverage(d+"/English/Fnt/", d+"/English/")
if err != nil {
  panic(err)
}

image, _ := gocr.ReadImage(d + "/imagetext_3.png")
s := gocr.NewCNNPredictorFromDir(modelPath + "sample1/model.cbor")

strings := gocr.ScanToStrings(s, image)
for _, s := range strings {
  fmt.Println(s)
}

```

# License
gocr is released under the Apache 2.0 License. se LICENSE for details.
