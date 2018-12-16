package view

import (
	"fmt"
	"github.com/andlabs/ui"
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
	s.Matrix[4] = x
	s.Matrix[5] = y
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
	Click   chan pixel.Vec
}

// Fonction lancée à l'initialisation du restaurant
// Boucle principale de l'interface graphique
func (w *Window) Draw() {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	fpsTxt := text.New(pixel.V(20, 20), atlas)
	fpsTxt.Color = image.Black
	frames := 0
	sec := time.Tick(time.Second)
	refresh := time.Tick(time.Second / time.Duration(60))
	for !w.Window.Closed() {
		if w.Window.JustPressed(pixelgl.MouseButtonLeft) {
			mp := w.Window.MousePosition()
			w.Click <- mp
		}
		w.Window.Clear(image.Black)
		for i := 0; i < len(w.Sprites); i++ {
			w.Sprites[i].PxlSprite.Draw(w.Window, w.Sprites[i].Matrix)
		}
		fpsTxt.Draw(w.Window, pixel.IM)
		select {
		case <-sec:
			fpsTxt.Clear()
			fmt.Fprintf(fpsTxt, "%v", frames)
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
		Fin:   make(chan bool),
		Click: make(chan pixel.Vec),
	}
	_ = win.NewSprite("ressources/map.png", 2)
	go win.Draw()
	return &win
}

// Crée une petite fenêtre avec du texte
// Utilisé pour décrire les éléments clickés
func Popup(title, content string) {
	ui.QueueMain(func() {
		win := ui.NewWindow(title, 300, 200, false)
		win.SetMargined(true)
		win.SetBorderless(true)
		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel(title+"\n\n"+content), false)
		button := ui.NewButton("Fermer")
		box.Append(button, false)
		win.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			win.Destroy()
		})
		win.Show()
	})
	// Pixel est complètement cassé avec plusieurs fenêtres
	/*
		w, err := pixelgl.NewWindow(pixelgl.WindowConfig{
			Title:  title,
			Bounds: pixel.R(0, 0, 300, 200),
		})
		if err != nil {
			panic(err)
		}
		atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		txt := text.New(pixel.V(180, 10), atlas)
		fmt.Fprintf(txt, content)
		fmt.Println(content)
		txt.Draw(w, pixel.IM)
		w.Update()
		for !w.Closed() {
		}
		w.Destroy()
	*/
}

// Vérifie si une entité est cliquée.
// Entrées: rectangle et matrice de l'entité, vecteur du curseur.
func CheckIfClicked(rect pixel.Rect, mat pixel.Matrix, vect pixel.Vec) bool {
	vect.X += (rect.Max.X / 2) * mat[0]
	vect.Y += (rect.Max.Y / 2) * mat[3]
	if (rect.Min.X+mat[4] < vect.X) && (rect.Max.X*mat[0]+mat[4] > vect.X) {
		if (rect.Min.Y+mat[5] < vect.Y) && (rect.Max.Y*mat[3]+mat[5] > vect.Y) {
			return true
		}
	}
	return false
}
