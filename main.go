// Partie "vue" du projet C#/.NET (sans .NET ni C#)
// Affiche et gère le restaurant, communique avec un serveur qui gère la cuisine.
// Implémente la vue et la partie du contrôleur gérant la salle du resto
// Le serveur gère le contrôleur de la cuisine et le modèle

package main

import (
	"fmt"
	"time"

	//"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/JamesMcAvoy/resto/src/controller"
	//"github.com/JamesMcAvoy/resto/src/view"
)

const (
	width  = 1280
	height = 720
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
	fmt.Println(game)
}

func main() {
	pixelgl.Run(run)
}
