package controller

import (
	"github.com/JamesMcAvoy/resto/src/view"
)

type Personne interface {
	Act()
}

type Client struct {
	Sprite *view.Sprite
}

func NewClient(r *Resto) Client {
	var c Client
	c.Sprite = r.Win.NewSprite("ressources/LeStig.png", 0.2)
	return c
}

func (c Client) Act() {
	c.Sprite.Move(2, 1)
}

type MaitreHotel struct {
	Sprite *view.Sprite
}

func NewMaitreHotel(r *Resto) MaitreHotel {
	var m MaitreHotel
	m.Sprite = r.Win.NewSprite("ressources/maitrehotel.png", 1)
	m.Sprite.Pos(40, 550)
	return m
}

func (m MaitreHotel) Act() {

}
