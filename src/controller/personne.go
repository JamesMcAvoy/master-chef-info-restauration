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

// CLIENT

// Client
type Client struct {
	Resto     *Resto
	Nom       string
	Sprite    *view.Sprite
	Etat      string
	Restant   int
	EstArrivé bool
}

// Constructeur de clients
func NewClient(r *Resto) *Client {
	var c Client
	c.Resto = r
	c.Nom = "Client"
	c.Sprite = r.Win.NewSprite("ressources/LeStig.png", 0.2)
	c.Sprite.Pos(20, 350)
	c.Etat = "Se demande si le restaurant est ouvert"
	c.Restant = 30
	r.Clickables = append(r.Clickables, &c)
	r.Personnes = append(r.Personnes, &c)
	return &c
}

// Ouvre le popup décrivant l'état du client quand il est cliqué
func (c *Client) CheckClick(mousePos pixel.Vec) bool {
	if view.CheckIfClicked(c.Sprite.PxlSprite.Picture().Bounds(), c.Sprite.Matrix, mousePos) {
		go view.Popup(c.Nom, c.String())
		return true
	}
	return false
}

// Stringer du client, sera affiché dans le popup décrivant le client
func (c *Client) String() string {
	return fmt.Sprintf("%s", c.Etat)
}

// Action effectuée par le client à chaque tick du restaurant
func (c *Client) Act() {
	c.Restant--
	switch c.Etat {
	case "S'en va":
		c.Sprite.Move(-2, 0)
	case "Se dirige vers le maître d'hôtel":
		if c.Resto.MaitreHotel.Sprite.Matrix[4]+50 > c.Sprite.Matrix[4] {
			c.Sprite.Move(2, 0)
		} else {
			c.EstArrivé = true
		}
		if c.Resto.MaitreHotel.Sprite.Matrix[5] > c.Sprite.Matrix[5] {
			c.Sprite.Move(0, 2)
		} else {
			if c.EstArrivé {
				c.Restant = 0
			}
		}
	}
	if c.Restant == 0 {
		switch c.Etat {
		case "Se demande si le restaurant est ouvert":
			c.Restant = 30
			if c.Resto.EstOuvert() {
				c.Etat = "Se dirige vers le maître d'hôtel"
			} else {
				c.Etat = "S'en va"
			}
		case "Se dirige vers le maitre d'hotel":
		}
	}
}

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
		m.ProchainClient = rand.Intn(300)
	}
}
