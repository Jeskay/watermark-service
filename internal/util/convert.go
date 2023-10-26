package util

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
)

func ByteToImage(image []byte, encoding string) image.Image {
	switch encoding {
	case ".png":
		img, err := png.Decode(bytes.NewReader(image))
		if err != nil {
			return nil
		}
		return img
	case ".jpg":
		img, err := jpeg.Decode(bytes.NewReader(image))
		if err != nil {
			return nil
		}
		return img
	default:
		return nil
	}
}

func ImageToBytes(image image.Image, encoding string) []byte {
	var buffer bytes.Buffer
	switch encoding {
	case ".png":
		png.Encode(&buffer, image)
		return buffer.Bytes()
	case ".jpg":
		jpeg.Encode(&buffer, image, nil)
		return buffer.Bytes()
	default:
		return nil
	}
}
