package controller

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"time"
)

// Struct restaurant
type Resto struct {
	Window   *pixelgl.Window
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
func NewResto(win *pixelgl.Window, temps, accel int, pause bool, h [][2]float64, e, p, d []string) *Resto {
	resto := Resto{
		Window:   win,
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
	return &resto
}

// Incrémente le temps dans le restaurant
func (r *Resto) incTick() {
	for {
		fmt.Println(r.Temps)
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
