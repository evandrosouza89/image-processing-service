package main

import (
	"fmt"
	"github.com/h2non/bimg"
)

func resize(image []byte, parameters map[string]interface{}) ([]byte, error) {

	width := int(parameters["width"].(float64))

	if width < 1 || width > MaxResizePixels {
		return nil, fmt.Errorf("invalid width: %d, must be greater than 1 and less than %d", width, MaxResizePixels)
	}

	height := int(parameters["height"].(float64))

	if height < 1 || height > MaxResizePixels {
		return nil, fmt.Errorf("invalid height: %d, must be greater than 1 and less than %d", height, MaxResizePixels)
	}

	newImage, err := bimg.NewImage(image).ForceResize(width, height)
	if err != nil {
		return nil, err
	}

	return newImage, nil

}

func enlarge(image []byte, parameters map[string]interface{}) ([]byte, error) {

	img := bimg.NewImage(image)

	size, err := img.Size()
	if err != nil {
		return nil, err
	}

	percentage := int(parameters["percentage"].(float64))

	if percentage == 0 || percentage == 100 || percentage < -MinEnlargePercentage || percentage > MaxEnlargePercentage {
		return nil, fmt.Errorf("invalid percentage: %d, must be different than 0, different than 100, greater than %d and less than %d", percentage, MinEnlargePercentage, MaxEnlargePercentage)
	}

	newImage, err := img.Enlarge(size.Width*percentage/100, size.Height*percentage/100)
	if err != nil {
		return nil, err
	}

	return newImage, nil

}

func rotate(image []byte, parameters map[string]interface{}) ([]byte, error) {

	angle := int(parameters["angle"].(float64))

	var rotateAngle bimg.Angle

	switch angle {
	case 90:
		rotateAngle = bimg.D90
	case 180:
		rotateAngle = bimg.D180
	case 270:
		rotateAngle = bimg.D270
	default:
		return nil, fmt.Errorf("invalid rotation angle: %d, must be 90, 180, or 270", angle)
	}

	newImage, err := bimg.NewImage(image).Rotate(rotateAngle)
	if err != nil {
		return nil, err
	}

	return newImage, nil

}
