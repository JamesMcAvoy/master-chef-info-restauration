package view

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	//"github.com/faiface/pixel/pixelgl"
)

// Charge une image pour la librairie graphique
func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
