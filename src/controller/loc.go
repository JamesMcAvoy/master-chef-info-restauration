package controller

import (
	"flag"
	"time"

	"github.com/JamesMcAvoy/resto/src/view"
)

// Struct restaurant
type Resto struct {
	Win      *view.Window
	Temps    int
	accel    int
	tick     <-chan time.Time
	Pause    bool
	Horaires [][2]float64
	Entrees  []string
	Plats    []string
	Desserts []string
}

// Constructeur de restaurant
func NewResto(width, height, temps, accel, i int, pause bool, h [][2]float64, e, p, d []string) *Resto {
	// La librairie Pixel DOIT utiliser le thread principal, elle ne peut pas s'exécuter ailleurs.
	// Mais le thread principal n'est pas disponible pendant les tests unitaires.
	// Les tests interagissant avec la librairie graphique sont impossibles, d'où cette condition.
	var win *view.Window
	if flag.Lookup("test.v") == nil {
		win = view.NewWindow(width, height, i)
	} else {
		win = nil
	}
	resto := Resto{
		Win:      win,
		Temps:    temps,
		accel:    accel,
		Pause:    pause,
		Horaires: h,
		Entrees:  e,
		Plats:    p,
		Desserts: d,
	}
	resto.tick = time.Tick(time.Second / time.Duration(accel))
	go resto.incTick()
	go tmpFuncPourTesterLaffichage(&resto)
	return &resto
}

func tmpFuncPourTesterLaffichage(r *Resto) {
	sprite := r.Win.AddSprite("ressources/serveur.png")
	sprite.Matrix = sprite.Matrix.Scaled(r.Win.Window.Bounds().Center(), 0.5)
}

// Incrémente le temps dans le restaurant
func (r *Resto) incTick() {
	for {
		select {
		case <-r.tick:
			if r.Pause == false {
				if r.Temps < 86400 {
					r.Temps++
				} else {
					r.Temps = 0
				}
			}
		}
	}
}
