package view

import (
	"testing"
)

type PictureTest struct {
	path   string
	sprite bool
	err    bool
}

func TestLoadPicture(t *testing.T) {
	want := []PictureTest{
		{"../../ressources/map.png", true, false},
		{"img_qui_existe_pas.png", false, true},
	}
	for _, r := range want {
		// mff juste teste si l'erreur/le sprite n'est pas nil
		result, err := LoadPicture(r.path)
		if (result != nil) != r.sprite || (err != nil) != r.err {
			t.Errorf("LoadPicture(%q): pixel.Picture, error est nil: %t, %t", r.path, r.sprite, r.err)
		}
	}
}
