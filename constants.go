package main

const (
	MaxBatchSize         = 3
	MaxUploadSize        = 350 << 10 // 350 KB
	MaxResizePixels      = 1024
	MinEnlargePercentage = 1
	MaxEnlargePercentage = 200
)

var SupportedImageFormats = map[string]bool{
	".jpeg": true,
	".jpg":  true,
	".png":  true,
}

var SupportedOperations = map[string]func(file []byte, parameters map[string]interface{}) ([]byte, error){
	"resize":  resize,
	"enlarge": enlarge,
	"rotate":  rotate,
}
