package view

import (
	"github.com/faiface/pixel"
	"image"
	// Importé pour initialiser des png
	_ "image/png"
	"os"
)

// LoadPicture charge une image pour la librairie graphique
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
