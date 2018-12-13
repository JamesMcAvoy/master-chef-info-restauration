package controller

import (
	"time"

	"github.com/JamesMcAvoy/resto/src/view"
)

// Struct restaurant
type Resto struct {
	Win       *view.Window
	Temps     int
	accel     int
	tick      <-chan time.Time
	Pause     bool
	Horaires  [][2]float64
	Entrees   []string
	Plats     []string
	Desserts  []string
	Personnes []Personne
}

// Constructeur de restaurant
func NewResto(width, height, temps, accel, i int, pause bool, h [][2]float64, e, p, d []string) *Resto {
	var win *view.Window
	win = view.NewWindow(width, height, i)
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
	go resto.loop()
	resto.Personnes = append(resto.Personnes, NewMaitreHotel(&resto))
	resto.Personnes = append(resto.Personnes, NewClient(&resto))
	return &resto
}

// Incrémente le temps dans le restaurant
func (r *Resto) loop() {
	for {
		select {
		case <-r.tick:
			for _, p := range r.Personnes {
				p.Act()
			}
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
