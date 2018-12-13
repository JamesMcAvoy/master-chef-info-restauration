package controller

import (
	"fmt"
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
	c.Sprite = r.Win.NewSprite("ressources/LeStig.png")
	c.Sprite.Matrix = c.Sprite.Matrix.Scaled(r.Win.Window.Bounds().Center(), 0.1)
	return c
}

func (c Client) Act() {
	fmt.Println("bonjour je suis un client")
}
