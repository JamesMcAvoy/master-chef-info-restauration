package view

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"image"
	"time"
)

// Sprite: Struct contenant le sprite à afficher et sa matrice.
// Toutes les entités s'affichant à l'écran devront implémenter *Sprite.
// Le constructeur place également un pointeur vers l'objet créé dans la fenêtre.
// Donc il suffit de modifier Objet.Sprite.Matrix pour changer sa position à l'écran.
type Sprite struct {
	PxlSprite *pixel.Sprite
	Matrix    pixel.Matrix
}

// Déplace le sprite de x pixels horizontalement et y verticalement
func (s *Sprite) Move(x, y float64) {
	s.Matrix = s.Matrix.Moved(pixel.V(x, y))
}

func (s *Sprite) Pos(x, y float64) {
	s.Matrix = pixel.IM.Moved(pixel.V(x, y))
}

// Ajoute un sprite à l'interface graphique
func (w *Window) NewSprite(path string, scale float64) *Sprite {
	img, err := LoadPicture(path)
	if err != nil {
		panic(err)
	}
	var sprite Sprite
	sprite.PxlSprite = pixel.NewSprite(img, img.Bounds())
	sprite.Matrix = pixel.IM.Moved(w.Window.Bounds().Center())
	sprite.Matrix = sprite.Matrix.Scaled(w.Window.Bounds().Center(), scale)
	w.Sprites = append(w.Sprites, &sprite)
	return &sprite

}

// Window: Chaque restaurant possède une fenêtre. Chaque fenêtre possède un array de pointeurs de sprite.
type Window struct {
	Window  *pixelgl.Window
	Sprites []*Sprite
	Fin     chan bool
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
	_ = win.NewSprite("ressources/map.png", 2)
	go win.Draw()
	return &win
}
