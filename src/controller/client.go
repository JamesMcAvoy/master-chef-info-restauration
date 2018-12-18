package controller

import (
	"fmt"
	"github.com/JamesMcAvoy/resto/src/view"
	"github.com/faiface/pixel"
	"math/rand"
)

// Client représente un groupe de clients
type Client struct {
	Resto     *Resto
	Nom       string
	Sprite    *view.Sprite
	Etat      string
	Restant   int
	EstArrivé bool
	Taille    int
	Table     *Table
}

// NewClient est le onstructeur de clients
func NewClient(r *Resto) *Client {
	var c Client
	c.Resto = r
	c.Nom = "Clients"
	c.Sprite = r.Win.NewSprite("ressources/LeStig.png", 0.2)
	c.Sprite.Pos(20, 350)
	c.Etat = "Se demande si le restaurant est ouvert"
	c.Restant = 30
	// Puissance pour avoir une plus grande probabilité de petits groupes
	c.Taille = int(rand.Float64()*rand.Float64()*5+1) * 2
	r.Clickables = append(r.Clickables, &c)
	r.Personnes = append(r.Personnes, &c)
	return &c
}

// CheckClick ouvre le popup décrivant l'état du client quand il est cliqué
func (c *Client) CheckClick(mousePos pixel.Vec) bool {
	if view.CheckIfClicked(c.Sprite.PxlSprite.Picture().Bounds(), c.Sprite.Matrix, mousePos) {
		go view.Popup(c.Nom, c.String())
		return true
	}
	return false
}

// Stringer du client, sera affiché dans le popup décrivant le client
func (c *Client) String() string {
	return fmt.Sprintf("Groupe de %v personnes\n\n%s", c.Taille, c.Etat)
}

// Act est l'action effectuée par le client à chaque tick du restaurant
func (c *Client) Act() {
	if c.Restant > 0 {
		c.Restant--
	}
	switch c.Etat {
	case "S'en va":
		c.Sprite.Move(-2, 0)
	case "Se dirige vers le maître d'hôtel":
		if c.Sprite.Goto(c.Resto.MaitreHotel.Sprite, 30, 0) {
			c.Restant = 0
		}
	case "Se dirige vers une table":
		if c.Sprite.Goto(c.Table.Sprite, 0, 50) {
			c.Etat = "Choisit un plat"
			c.Restant = 150
		}
	}
	if c.Restant == 0 {
		switch c.Etat {
		case "Se demande si le restaurant est ouvert":
			c.Restant = -1
			if c.Resto.EstOuvert() {
				c.Etat = "Se dirige vers le maître d'hôtel"
			} else {
				c.Etat = "S'en va"
			}
		case "Se dirige vers le maître d'hôtel":
			c.Resto.MaitreHotel.Queue = append(c.Resto.MaitreHotel.Queue, c)
			c.Etat = "En attente d'attribution de table"
		case "Choisit un plat":
			//c.Table.Carré.AppelServeur
		}
	}
}
