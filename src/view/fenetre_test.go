package view

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"os"
	"testing"
)

// Donne le thread principal à pixel
func TestMain(m *testing.M) {
	err := os.Chdir(os.Getenv("GOPATH") + "/src/github.com/JamesMcAvoy/resto")
	if err != nil {
		panic(err)
	}
	pixelgl.Run(func() {
		os.Exit(m.Run())
	})
}

func TestPosMove(t *testing.T) {
	win := NewWindow(1280, 704)
	cases := []struct {
		sprite                   *Sprite
		posX, posY, moveX, moveY float64
		expectedX, expectedY     float64
	}{
		{win.NewSprite("ressources/plat.png", 2), 0, 0, 0, 0, 0, 0},
		{win.NewSprite("ressources/plat.png", 2), 0, 0, 30, 50, 30, 50},
		{win.NewSprite("ressources/plat.png", 2), 30, 50, 0, 0, 30, 50},
		{win.NewSprite("ressources/plat.png", 2), -30, 50, 75, -15, 45, 35},
		{win.NewSprite("ressources/plat.png", 2), -30, 50, -75, -15, -105, 35},
	}
	for _, v := range cases {
		v.sprite.Pos(v.posX, v.posY)
		v.sprite.Move(v.moveX, v.moveY)
		if v.sprite.Matrix[4] != v.expectedX || v.sprite.Matrix[5] != v.expectedY {
			t.Errorf("sprite.Pos(%v, %v).Move(%v, %v): attendu %v, %v, reçu %v, %v",
				v.posX, v.posY, v.moveX, v.moveY, v.expectedX, v.expectedY, v.sprite.Matrix[4], v.sprite.Matrix[5])
		}
	}
}

var CheckIfClickedTest = []struct {
	// Rectangle dans lequel l'image est affichée,
	// matrice appliquée pour le transformer/déplacer,
	// vecteur représentant la position du curseur
	rect     pixel.Rect
	mat      pixel.Matrix
	vect     pixel.Vec
	expected bool
}{
	// Image de 0x0, ne doit pas être cliquée
	{pixel.R(0, 0, 0, 0), pixel.IM, pixel.V(0, 0), false},
	// Image de 0x50, placée à la position 50, 50, ne doit pas être cliquée à 76, 76
	// La transformation de la matrice est appliquée au centre du rectangle, pas au bord
	{pixel.R(0, 0, 50, 50), pixel.IM.Moved(pixel.V(50, 50)), pixel.V(76, 76), false},
	{pixel.R(0, 0, 50, 50), pixel.IM.Moved(pixel.V(0, 50)), pixel.V(75, 175), false},
	{pixel.R(0, 0, 50, 50), pixel.IM.Moved(pixel.V(50, 50)), pixel.V(25, 25), false},

	// Image de 50x50, placée à la position 50, 0, taille x2, ne doit pas être cliquée à 130, 75
	{pixel.R(0, 0, 50, 50), pixel.Matrix{2, 0, 0, 2, 50, 0}, pixel.V(130, 75), false},
	// Image de 50x50, placée à la position 0, 50, taille x0.5, ne doit pas être cliquée à 80, 30
	{pixel.R(0, 0, 50, 50), pixel.Matrix{0.5, 0, 0, 0.5, 0, 50}, pixel.V(30, 80), false},
	{pixel.R(0, 0, 50, 50), pixel.Matrix{0.5, 0, 20, 0.5, 0, 20}, pixel.V(10, 10), true},
}

func TestCheckIfClicked(t *testing.T) {
	for _, v := range CheckIfClickedTest {
		output := CheckIfClicked(v.rect, v.mat, v.vect)
		if output != v.expected {
			t.Errorf("CheckIfClicked(%s, %s, %s): %t attendu, %t reçu",
				v.rect, v.mat, v.vect, v.expected, output)
		}
	}
}
