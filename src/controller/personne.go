package controller

import (
	"fmt"
	"github.com/JamesMcAvoy/resto/src/view"
	"github.com/faiface/pixel"
)

// interface Clickable: Tous les objets sur lesquels ont peut cliquer l'implémentent.
// Cliquer dessus fait apparaître un popup décrivant l'objet.
type Clickable interface {
	CheckClick(pixel.Vec) bool
}

// Interface Personne: Le restaurant exécute la méthode Act de tous les
// humains dans le restaurant à chaque tick.
type Personne interface {
	Act()
	CheckClick(pixel.Vec) bool
}

// CLIENT

// Client
type Client struct {
	Nom    string
	Sprite *view.Sprite
	Etat   string
}

// Constructeur de clients
func NewClient(r *Resto) Client {
	var c Client
	c.Nom = "Client"
	c.Sprite = r.Win.NewSprite("ressources/LeStig.png", 0.2)
	c.Sprite.Pos(10, 40)
	c.Etat = "En attente d'attribution de table"
	return c
}

// Ouvre le popup décrivant l'état du client quand il est cliqué
func (c Client) CheckClick(mousePos pixel.Vec) bool {
	if view.CheckIfClicked(c.Sprite.PxlSprite.Picture().Bounds(), c.Sprite.Matrix, mousePos) {
		go view.Popup(c.Nom, c.String())
		return true
	}
	return false
}

// Stringer du client, sera affiché dans le popup décrivant le client
func (c Client) String() string {
	return fmt.Sprintf("%s", c.Etat)
}

// Action effectuée par le client à chaque tick du restaurant
func (c Client) Act() {
	c.Sprite.Move(2, 1)
}

// MAITRE D'HOTEL

// Maître d'hôtel
type MaitreHotel struct {
	Nom    string
	Sprite *view.Sprite
}

// Constructeur de maître d'ĥôtel
func NewMaitreHotel(r *Resto) MaitreHotel {
	var m MaitreHotel
	m.Nom = "Maître d'ĥôtel"
	m.Sprite = r.Win.NewSprite("ressources/maitrehotel.png", 1)
	m.Sprite.Pos(140, 550)
	return m
}

// Ouvre le popup décrivant l'état du maître d'ĥôtel quand il est cliqué
func (m MaitreHotel) CheckClick(mousePos pixel.Vec) bool {
	if view.CheckIfClicked(m.Sprite.PxlSprite.Picture().Bounds(), m.Sprite.Matrix, mousePos) {
		go view.Popup(m.Nom, m.String())
		return true
	}
	return false
}

// Stringer du maître d'hôtel, sera affiché dans le popup le décrivant
func (m MaitreHotel) String() string {
	return fmt.Sprintf("Maître d'hôtel qui attend")
}

// Action effectuée par le maître d'hôtel à chaque tick du restaurant
func (m MaitreHotel) Act() {

}
