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

// Serveur
type Serveur struct {
	Carré  *Carré
	Nom    string
	Etat   string
	Sprite *view.Sprite
	Client *Client // Le client dont le serveur s'occupe
}

// Constructeur de serveur
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

// Ouvre le popup décrivant l'état du serveur quand il est cliqué
func (s *Serveur) CheckClick(mousePos pixel.Vec) bool {
	if view.CheckIfClicked(s.Sprite.PxlSprite.Picture().Bounds(), s.Sprite.Matrix, mousePos) {
		go view.Popup(s.Nom, s.String())
		return true
	}
	return false
}

func (s *Serveur) String() string {
	return s.Etat
}

// MAITRE D'HOTEL

// Maître d'hôtel
type MaitreHotel struct {
	Resto          *Resto
	Nom            string
	Etat           string
	Sprite         *view.Sprite
	Queue          []*Client
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
	r.Personnes = append(r.Personnes, &m)
	r.Clickables = append(r.Clickables, &m)
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
		if table := m.TableLibre(client.Taille); table != nil {
			for _, serveur := range table.Carré.Serveurs {
				if serveur.Etat == "Libre" {
					serveur.Etat = "Se dirige vers un client pour le placer"
					serveur.Client = client
					client.Table = table
					m.Queue = append(m.Queue[:i], m.Queue[i+1:]...)
					break
				}
			}
		}
	}
}

// Retourne la table libre la plus petite pour le groupe qui arrive,
// retourne nil si pas de table disponible
func (m *MaitreHotel) TableLibre(taille int) *Table {
	for taille <= 10 {
		for i := range m.Resto.Carrés {
			for j := range m.Resto.Carrés[i].Tables {
				if m.Resto.Carrés[i].Tables[j].Taille == taille && !m.Resto.Carrés[i].Tables[j].Occupée {
					m.Resto.Carrés[i].Tables[j].Occupée = true
					return m.Resto.Carrés[i].Tables[j]
				}
			}
		}
		taille += 1
	}
	return nil
}
