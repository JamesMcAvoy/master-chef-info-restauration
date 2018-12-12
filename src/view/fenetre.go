package view

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"image"
	"time"
	//"unicode"
)

// Sprite: Struct contenant le sprite à afficher et sa matrice.
type Sprite struct {
	PxlSprite *pixel.Sprite
	Matrix    pixel.Matrix
}

// Window: Chaque restaurant possède une fenêtre. Chaque fenêtre possède un array de pointeurs de sprite.
// Toutes les entités (serveurs, tables...) affichés à l'écran devront avoir un pointeur de sprite
// pointant vers un élément de ce tableau.
// Ainsi, simplement modifier la matrice de l'élément dans le controller changera automatiquement sa
// position dans l'interface graphique.
type Window struct {
	Window  *pixelgl.Window
	Sprites []*Sprite
	Fin     chan bool
}

// Ajoute un sprite à l'interface graphique
func (w *Window) AddSprite(path string) *Sprite {
	img, err := LoadPicture(path)
	if err != nil {
		panic(err)
	}
	var sprite Sprite
	sprite.PxlSprite = pixel.NewSprite(img, img.Bounds())
	sprite.Matrix = pixel.IM.Moved(w.Window.Bounds().Center())
	sprite.Matrix = sprite.Matrix.Scaled(w.Window.Bounds().Center(), 2)
	w.Sprites = append(w.Sprites, &sprite)
	return &sprite

}

// Fonction lancée à l'initialisation du restaurant
// Boucle principale de l'interface graphique
func (w *Window) Draw() {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	fpsTxt := text.New(pixel.V(20, 20), atlas)
	frames := 0
	sec := time.Tick(time.Second)
	refresh := time.Tick(time.Second / time.Duration(60))
	for !w.Window.Closed() {
		w.Window.Clear(image.Black)
		fmt.Println(w.Sprites)
		for i := 0; i < len(w.Sprites); i++ {
			w.Sprites[i].PxlSprite.Draw(w.Window, w.Sprites[i].Matrix)
		}
		fpsTxt.Draw(w.Window, pixel.IM)
		select {
		case <-sec:
			fpsTxt.Clear()
			fmt.Fprintf(fpsTxt, "FPS: %v", frames)
			frames = 0
		default:
			<-refresh
		}
		frames++
		w.Window.Update()
	}
	w.Fin <- true
	w.Window.Destroy()
}

// Crée une fenêtre
func NewWindow(width, height, i int) *Window {
	w, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  fmt.Sprintf("La salle du resto %v oui", i),
		Bounds: pixel.R(0, 0, float64(width), float64(height)),
	})
	if err != nil {
		panic(err)
	}
	win := Window{
		Window: w,
		//Sprites: []*Sprite{sprite},
		Fin: make(chan bool),
	}
	_ = win.AddSprite("ressources/map.png")
	go win.Draw()
	return &win
}
