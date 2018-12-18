package controller

import (
	"fmt"
	"github.com/JamesMcAvoy/resto/src/view"
	"github.com/faiface/pixel"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// Resto représente un restaurant
type Resto struct {
	Win         *view.Window
	Temps       int
	accel       float64
	tick        <-chan time.Time
	MaitreHotel *MaitreHotel
	Horaires    [][2]float64
	Entrees     []string
	Plats       []string
	Desserts    []string
	Carrés      []*Carré
	Clickables  []Clickable
	Personnes   []Personne
}

// NewResto crée un restaurant, place ses tables et ses serveurs
func NewResto(width, height, temps, accel int, h [][2]float64, e, p, d []string, carrés []interface{}) *Resto {
	// Crée la fenêtre
	var win *view.Window
	win = view.NewWindow(width, height)
	resto := Resto{
		Win:      win,
		Temps:    temps,
		accel:    float64(accel),
		Horaires: h,
		Entrees:  e,
		Plats:    p,
		Desserts: d,
	}
	for i := range h {
		for j := range h[i] {
			h[i][j] = h[i][j] * 3600
		}
	}
	resto.Horaires = h
	// Crée les carrés
	carPos := Répartit(width, height, len(carrés))
	for i := range carPos {
		if carPos[i][0] == 0 { // Empêche les tables d'apparaître dans les murs
			carPos[i][0] += 80
		}
		if carPos[i][1] == 0 {
			carPos[i][1] += 30
		}
		if carPos[i][2] == width {
			carPos[i][2] -= 30
		}
		if carPos[i][3] == height {
			carPos[i][3] -= 40
		}
		resto.Carrés = append(resto.Carrés, NewCarré(carPos[i], carrés[i].(map[string]interface{}), &resto))
		// Place les tables
	}
	// Crée le maître d'hôtel et lance le restaurant
	resto.tick = time.Tick(time.Second / time.Duration(accel))
	resto.MaitreHotel = NewMaitreHotel(&resto)
	go resto.loop()
	return &resto
}

// EstOuvert vérifie si le restaurant est ouvert
func (r *Resto) EstOuvert() bool {
	for _, v := range r.Horaires {
		if float64(r.Temps) > v[0] && float64(r.Temps) < v[1] {
			return true
		}
	}
	return false
}

// loop est la boucle principale du restaurant
func (r *Resto) loop() {
	min := 0
	prev := 0
	for {
		min = r.Temps % 3600 / 60
		if min != prev {
			if min < 10 {
				r.Win.Window.SetTitle(fmt.Sprintf("La salle du resto, %vh0%v", r.Temps/3600, min))
			} else {
				r.Win.Window.SetTitle(fmt.Sprintf("La salle du resto, %vh%v", r.Temps/3600, min))
			}
			prev = min
		}
		select {
		case mousePos := <-r.Win.Click:
			for i := range r.Clickables {
				if r.Clickables[i].CheckClick(mousePos) {
					break
				}
			}
		case scroll := <-r.Win.Scroll:
			acc := r.accel + scroll/10*r.accel
			if acc > 1.2 {
				r.accel = acc
				r.tick = time.Tick(time.Second / time.Duration(r.accel))
			} else {
				if scroll > 0 {
					r.accel = 1.1
					r.tick = time.Tick(time.Second / time.Duration(1))
				} else {
					r.tick = nil
				}
			}
		case <-r.tick:
			for _, p := range r.Personnes {
				p.Act()
			}
			if r.Temps < 86400 {
				r.Temps++
			} else {
				r.Temps = 0
			}
		}
	}
}

// Répartit répartit un espace de dimensions width, height en nb carrés de façon à ce que le nombre
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

// repLoop: juste une petite fonction pour ne pas répéter de code dans Répartit()
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

// Carré est un ensemble de tables dont un groupe de serveurs s'occupe
type Carré struct {
	// basGaucheX, basGaucheY, hautDroiteX, hautDroiteY
	Coords   [4]int
	Tables   []*Table
	Serveurs []*Serveur
	Resto    *Resto
}

// NewCarré crée un carré, lui attribue les tables et les serveurs
func NewCarré(pos [4]int, car map[string]interface{}, resto *Resto) *Carré {
	c := Carré{Coords: pos, Resto: resto}
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
	tablePos := Répartit(
		c.Coords[2]-c.Coords[0],
		c.Coords[3]-c.Coords[1],
		int(tableCount),
	)
	index := 0
	for j := range tablePos {
		go func(j, index int) {
			table := NewTable(
				tailles[index],
				[4]int{
					pos[0] + tablePos[j][0], pos[1] + tablePos[j][1],
					pos[0] + tablePos[j][2], pos[1] + tablePos[j][3]},
				&c,
			)
			c.Tables = append(c.Tables, table)
			resto.Clickables = append(resto.Clickables, table)
		}(j, index)
		index++
	}
	for i := 0.0; i <= tableCount/5; i++ {
		c.Serveurs = append(c.Serveurs, NewServeur(&c))
	}
	return &c
}

// Table représente une table
type Table struct {
	Sprite  *view.Sprite
	Carré   *Carré
	Taille  int
	Nom     string
	Coords  [4]int
	Occupée bool
}

// NewTable crée une table
func NewTable(taille int, coords [4]int, c *Carré) *Table {
	var t Table
	t.Nom = "Une table"
	t.Taille = taille
	t.Carré = c
	t.Sprite = c.Resto.Win.NewSprite(fmt.Sprintf("ressources/table%v.png", taille), 0.85)
	t.Sprite.Pos(float64((coords[2]+coords[0])/2), float64((coords[3]+coords[1])/2))
	time.Sleep(time.Second)
	return &t
}

// CheckClick ouvre un popup quand la table est cliquée
func (t *Table) CheckClick(mousePos pixel.Vec) bool {
	if view.CheckIfClicked(t.Sprite.PxlSprite.Picture().Bounds(), t.Sprite.Matrix, mousePos) {
		go view.Popup(t.Nom, t.String())
		return true
	}
	return false
}

func (t *Table) String() string {
	if t.Occupée {
		return fmt.Sprintf("Table de %v personnes occupée", t.Taille)
	}
	return fmt.Sprintf("Table de %v personnes libre", t.Taille)
}
