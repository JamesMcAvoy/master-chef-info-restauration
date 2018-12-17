package controller

import (
	"fmt"
	"github.com/JamesMcAvoy/resto/src/view"
	"github.com/faiface/pixel"
	"math/rand"
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

// SERVEUR

type Serveur struct{}

// MAITRE D'HOTEL

// Maître d'hôtel
type MaitreHotel struct {
	Resto          *Resto
	Nom            string
	Etat           string
	Sprite         *view.Sprite
	Queue          []Client
	ProchainClient int
}

// Constructeur de maître d'ĥôtel
func NewMaitreHotel(r *Resto) *MaitreHotel {
	var m MaitreHotel
	m.Nom = "Maître d'hôtel"
	m.Sprite = r.Win.NewSprite("ressources/maitrehotel.png", 1)
	m.Sprite.Pos(40, 550)
	m.ProchainClient = rand.Intn(300)
	m.Resto = r
	return &m
}

// Ouvre le popup décrivant l'état du maître d'ĥôtel quand il est cliqué
func (m *MaitreHotel) CheckClick(mousePos pixel.Vec) bool {
	if view.CheckIfClicked(m.Sprite.PxlSprite.Picture().Bounds(), m.Sprite.Matrix, mousePos) {
		go view.Popup(m.Nom, m.String())
		return true
	}
	return false
}

// Stringer du maître d'hôtel, sera affiché dans le popup le décrivant
func (m *MaitreHotel) String() string {
	return fmt.Sprintf("Temps avant l'arrivée du prochain client: %v", m.ProchainClient)
}

// Action effectuée par le maître d'hôtel à chaque tick du restaurant.
// Arrivée des clients
func (m *MaitreHotel) Act() {
	m.ProchainClient--
	// Arrivée des clients
	if m.ProchainClient == 0 {
		m.Queue = append(m.Queue, *NewClient(m.Resto))
		m.ProchainClient = rand.Intn(20) + 1
	}
}
