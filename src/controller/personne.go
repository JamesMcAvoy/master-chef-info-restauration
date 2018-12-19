package controller

import (
	"fmt"
	"github.com/JamesMcAvoy/resto/src/view"
	"github.com/faiface/pixel"
	"math/rand"
)

// Clickable est l'interface implémentée par tous les objets sur lesquels ont peut cliquer.
// Cliquer dessus fait apparaître un popup décrivant l'objet.
type Clickable interface {
	CheckClick(pixel.Vec) bool
}

// Personne est l'interface implémentée par tous les humains dans le restaurant.
// Leur méthode "Act" est appelée à chaque tick.
type Personne interface {
	Act()
	CheckClick(pixel.Vec) bool
}

// SERVEUR

// Serveur représente un serveur
type Serveur struct {
	Carré  *Carré
	Nom    string
	Etat   string
	Sprite *view.Sprite
	Client *Client // Le client dont le serveur s'occupe
}

// NewServeur construit un serveur
func NewServeur(c *Carré) *Serveur {
	var s Serveur
	s.Nom = "Un serveur"
	s.Sprite = c.Resto.Win.NewSprite("ressources/serveur.png", 0.3)
	s.Sprite.Pos(rand.Float64()*1000, rand.Float64()*700)
	s.Etat = "Libre"
	s.Carré = c
	c.Resto.Personnes = append(c.Resto.Personnes, &s)
	c.Resto.Clickables = append(c.Resto.Clickables, &s)
	return &s
}

// Act est la fonction exécutée par le serveur à chaque tick du restaurant
func (s *Serveur) Act() {
	switch s.Etat {
	case "Se dirige vers un client pour le placer":
		if s.Sprite.Goto(s.Client.Sprite, 30, 0) {
			s.Client.Etat = "Se dirige vers une table"
			s.Etat = "Se dirige vers une table pour placer les clients"
		}
	case "Se dirige vers une table pour placer les clients":
		if s.Sprite.Goto(s.Client.Table.Sprite, 40, 0) {
			s.Etat = "Libre"
		}
	}
}

// CheckClick ouvre le popup décrivant l'état du serveur quand il est cliqué
func (s *Serveur) CheckClick(mousePos pixel.Vec) bool {
	if view.CheckIfClicked(s.Sprite.PxlSprite.Frame(), s.Sprite.Matrix, mousePos) {
		go view.Popup(s.Nom, s.String())
		return true
	}
	return false
}

func (s *Serveur) String() string {
	return s.Etat
}

// MAITRE D'HOTEL

// MaitreHotel représente le Maître d'hôtel
type MaitreHotel struct {
	Resto          *Resto
	Nom            string
	Etat           string
	Sprite         *view.Sprite
	Queue          []*Client
	ProchainClient int
}

// NewMaitreHotel construit un maître d'ĥôtel
func NewMaitreHotel(r *Resto) *MaitreHotel {
	var m MaitreHotel
	m.Nom = "Maître d'hôtel"
	m.Sprite = r.Win.NewSprite("ressources/maitrehotel.png", 1)
	m.Sprite.Pos(40, 550)
	m.ProchainClient = rand.Intn(300)
	m.Resto = r
	r.Personnes = append(r.Personnes, &m)
	r.Clickables = append(r.Clickables, &m)
	return &m
}

// CheckClick ouvre le popup décrivant l'état du maître d'ĥôtel quand il est cliqué
func (m *MaitreHotel) CheckClick(mousePos pixel.Vec) bool {
	if view.CheckIfClicked(m.Sprite.PxlSprite.Frame(), m.Sprite.Matrix, mousePos) {
		go view.Popup(m.Nom, m.String())
		return true
	}
	return false
}

// Stringer du maître d'hôtel, sera affiché dans le popup le décrivant
func (m *MaitreHotel) String() string {
	return fmt.Sprintf("Temps avant l'arrivée du prochain client: %v", m.ProchainClient)
}

// Act est la fonction exécutée par le maître d'hôtel à chaque tick
func (m *MaitreHotel) Act() {
	m.ProchainClient--
	// Arrivée des clients
	if m.ProchainClient == 0 {
		NewClient(m.Resto)
		m.ProchainClient = rand.Intn(300) + 1
	}
	// Attribution d'une table à un client de la file, appel d'un serveur pour le placer
	//if len(m.Queue) > 0 {
	for i, client := range m.Queue {
		if client.Table == nil {
			m.AttribueTable(client)
		}
		if client.Table != nil {
			if serveur := client.Table.Carré.ServeurLibre(); serveur != nil {
				serveur.Etat = "Se dirige vers un client pour le placer"
				serveur.Client = client
				m.Queue = append(m.Queue[:i], m.Queue[i+1:]...)
			}
		}
	}
}

// TableLibre attribue la table libre la plus petite pour un groupe de client.
// Ne fait rien si aucune table n'est disponible.
func (m *MaitreHotel) AttribueTable(client *Client) {
	for taille := client.Taille; taille <= 10; taille++ {
		for _, carré := range m.Resto.Carrés {
			for _, table := range carré.Tables {
				if table.Taille == taille && !table.Occupée {
					client.Table = table
					table.Occupée = true
					return
				}
			}
		}
	}
}
