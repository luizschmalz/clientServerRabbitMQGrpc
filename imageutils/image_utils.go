package imageutils

import (
	"bytes"
	"image"
	"github.com/disintegration/imaging"
)

func ToGray(img image.Image) image.Image {
	return imaging.Grayscale(img)
}

func ImageToBytes(img image.Image) []byte {
	var buf bytes.Buffer
	_ = imaging.Encode(&buf, img, imaging.PNG)
	return buf.Bytes()
}

func BytesToImage(data []byte) (image.Image, error) {
	img, err := imaging.Decode(bytes.NewReader(data))
	return img, err
}
