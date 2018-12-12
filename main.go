// Partie "salle" du projet C#/.NET (sans .NET ni C#)
// Affiche et gère le restaurant, communique avec un serveur qui gère la cuisine.

package main

import (
	"github.com/faiface/pixel/pixelgl"
	"time"

	"github.com/JamesMcAvoy/resto/src/controller"
)

const (
	width  = 1280
	height = 704
	url    = "http://127.0.0.1:9090/"
	// Jusqu'à ce qu'on lie les 2 projets:
	acceleration = 60 // Accélération initiale du temps
	port         = 9090
)

func run() {
	// Jusqu'à ce qu'on lie les 2 projets:
	go Serv(port, acceleration)
	time.Sleep(50 * time.Millisecond)

	game := controller.NewGame(width, height, url)
	fin := make(chan bool)
	for i, r := range game.Restos {
		go func(i int, r *controller.Resto) {
			<-r.Win.Fin
			if len(game.Restos) > 1 {
				// Supprime le restaurant en évitant les memory leaks
				game.Restos[i] = game.Restos[len(game.Restos)-1]
				game.Restos[len(game.Restos)-1] = nil
				game.Restos = game.Restos[:len(game.Restos)-1]
			} else {
				fin <- true
			}
		}(i, r)
	}
	<-fin
}

func main() {
	pixelgl.Run(run)
}
