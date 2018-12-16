package controller

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/JamesMcAvoy/resto/src/view"
)

// Struct restaurant
type Resto struct {
	Win        *view.Window
	Temps      int
	accel      int
	tick       <-chan time.Time
	Pause      bool
	Horaires   [][2]float64
	Entrees    []string
	Plats      []string
	Desserts   []string
	Carrés     []Carré
	Clickables []Clickable
	Personnes  []Personne
}

// Constructeur de restaurant
func NewResto(width, height, temps, accel, i int, pause bool, h [][2]float64, e, p, d []string, carrés []interface{}) *Resto {
	// Crée la fenêtre
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
	// Crée les carrés
	carPos := Répartit(width, height, len(carrés))
	for i := range carPos {
		if carPos[i][0] == 0 { // Empêche les tables d'apparaître dans les murs
			carPos[i][0] += 50
		}
		//if carPos[i][1] == 0 {
		//	carPos[i][1] += 50
		//}
		if carPos[i][2] == width {
			carPos[i][2] -= 50
		}
		if carPos[i][3] == height {
			carPos[i][3] -= 50
		}
		resto.Carrés = append(resto.Carrés, Carré{Coords: carPos[i]})
		// Place les tables
		car := carrés[i].(map[string]interface{})
		var tailles []int
		tableCount := 0.0
		for k, v := range car {
			t, _ := strconv.Atoi(k)
			for i := 0.0; i < v.(float64); i++ {
				tailles = append(tailles, t)
			}
			tableCount += v.(float64)
		}
		for i := range tailles {
			j := rand.Intn(i + 1)
			tailles[i], tailles[j] = tailles[j], tailles[i]
		}
		tablePos := Répartit(resto.Carrés[i].Coords[2]-resto.Carrés[i].Coords[0], resto.Carrés[i].Coords[3]-resto.Carrés[i].Coords[1], int(tableCount))
		index := 0
		for j := range tablePos {
			go func(i, j, index int) {
				resto.Carrés[i].Tables = append(resto.Carrés[i].Tables, NewTable(
					tailles[index],
					[4]int{carPos[i][0] + tablePos[j][0], carPos[i][1] + tablePos[j][1], carPos[i][0] + tablePos[j][2], carPos[i][1] + tablePos[j][3]},
					&resto))
			}(i, j, index)
			index++
		}
	}
	// Crée le maître d'hôtel et lance le restaurant
	resto.tick = time.Tick(time.Second / time.Duration(accel))
	maitreHotel := NewMaitreHotel(&resto)
	resto.Personnes = append(resto.Personnes, maitreHotel)
	resto.Clickables = append(resto.Clickables, maitreHotel)
	go resto.loop()
	return &resto
}

// Incrémente le temps dans le restaurant
func (r *Resto) loop() {
	for {
		select {
		case mousePos := <-r.Win.Click:
			for i := len(r.Clickables) - 1; i >= 0; i-- {
				r.Clickables[i].CheckClick(mousePos)
			}
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

// Répartit un espace de dimensions width, height en nb carrés de façon à ce que le nombre
// de lignes et le nombre de colomnes soient le plus proche possible.
//
// Retourne un tableau de tableaux de coordonées des rectangles  sous la forme:
//
// [[XBasGauche, YBasGauche, XHautDroite, YHautDroite]]
//
// 1 objet: occupe tout l'espace,
// 2 objets: divisé en 2 verticalement,
// 3 objets: 1 colomne de 2 objets, 1 colomne de 1 objet
//
// 10 objets: 1 colomne de 4 objets, 2 colomnes de 3 objets
//
// 12 objets: 4 colomnes de 3 objets
//
// 45 objets: 3 colomnes de 7 objets, 3 colomnes de 6 objets
func Répartit(width, height, nb int) [][4]int {
	if nb == 0 {
		return nil
	}
	w := int(math.Ceil(math.Sqrt(float64(nb))))
	h := int(math.Sqrt(float64(nb)))
	returned := make([][4]int, nb)
	index := 0
	shift := 0
	switch {
	case nb < w*h:
		repLoop(width*(h-(w*h-nb))/h, height, h-(h*w-nb), w, &index, &shift, returned)
		repLoop(width*(w*h-nb)/h, height, h*w-nb, h, &index, &shift, returned)
	case nb == w*h:
		repLoop(width, height, w, h, &index, &shift, returned)
	case nb > w*h:
		repLoop(width*(nb-w*h)/w, height, nb-h*w, w, &index, &shift, returned)
		repLoop(width*(w-(nb-w*h))/w, height, w-(nb-h*w), h, &index, &shift, returned)
	}
	return returned
}

// Juste une petite fonction pour ne pas répéter de code dans Répartit()
func repLoop(width, height, w, h int, index, shift *int, returned [][4]int) {
	for i := 0; i < w; i++ {
		he := 0
		for j := 0; j < h; j++ {
			returned[*index] = [4]int{*shift, he, *shift + width/w, he + height/h}
			he += height / h
			*index++
		}
		*shift += width / w
	}
}

type Carré struct {
	// basGaucheX, basGaucheY, hautDroiteX, hautDroiteY
	Coords [4]int
	Tables []Table
}

type Table struct {
	Sprite *view.Sprite
	Taille int
	Nom    string
	Coords [4]int
}

func NewTable(taille int, coords [4]int, r *Resto) Table {
	var t Table
	t.Nom = "Une table"
	t.Taille = taille
	t.Sprite = r.Win.NewSprite(fmt.Sprintf("ressources/table%v.png", taille), 0.85)
	t.Sprite.Pos(float64((coords[2]+coords[0])/2), float64((coords[3]+coords[1])/2))
	return t
}
