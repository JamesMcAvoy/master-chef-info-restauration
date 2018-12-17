package controller

import (
	"fmt"
	"github.com/JamesMcAvoy/resto/src/view"
	"github.com/faiface/pixel"
	"math/rand"
)

// Client
type Client struct {
	Resto     *Resto
	Nom       string
	Sprite    *view.Sprite
	Etat      string
	Restant   int
	EstArrivé bool
	Taille    int
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
	// Puissance pour avoir une plus grande probabilité de petits groupes
	c.Taille = int(rand.Float64()*rand.Float64()*5+1) * 2
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
	return fmt.Sprintf("Groupe de %v personnes\n\n%s", c.Taille, c.Etat)
}

// Action effectuée par le client à chaque tick du restaurant
func (c *Client) Act() {
	c.Restant--
	switch c.Etat {
	case "S'en va":
		c.Sprite.Move(-2, 0)
	case "Se dirige vers le maître d'hôtel":
		c.Goto(c.Resto.MaitreHotel.Sprite)
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
			c.Etat = "En attente d'attribution de table"
		}
	}
}

func (c *Client) Goto(*view.Sprite) {
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
