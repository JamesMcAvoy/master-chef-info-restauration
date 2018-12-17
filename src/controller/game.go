package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Game représente l'"objet principal".
// Il contient les restos et la liaison avec le serveur
type Game struct {
	Url    string
	Restos []*Resto
}

// Initialisation des restos et connection au serveur
func NewGame(width, height int, url string) *Game {
	game := Game{
		Url: url,
	}
	bonjour := make(map[string]interface{})
	bonjour["type"] = "bonjour"
	initMap, err := game.Req(bonjour)
	if err != nil {
		panic(err)
	}
	tmp := initMap["restos"].([]interface{})
	m := make([]map[string]interface{}, len(tmp))
	for i, v := range tmp {
		m[i] = v.(map[string]interface{})
	}
	for _, r := range m {
		hor := r["horaires"].([]interface{})
		h := make([][2]float64, len(hor))
		for i, v := range hor {
			t := v.([]interface{})
			for j, val := range t {
				h[i][j] = val.(float64)
			}
		}
		en := r["entrees"].([]interface{})
		pl := r["plats"].([]interface{})
		de := r["desserts"].([]interface{})
		e := make([]string, len(en))
		p := make([]string, len(pl))
		d := make([]string, len(de))
		IntToStr(en, e)
		IntToStr(pl, p)
		IntToStr(de, d)
		game.Restos = append(game.Restos, NewResto(
			width, height, int(r["temps"].(float64)), int(r["acceleration"].(float64)),
			h, e, p, d, r["carres"].([]interface{})))
		// Pixel semble plus être un peu moins cassé quand il n'a pas à créer plusieurs fenêtres en même temps
		// (c'est peut-être totalement faux)
		time.Sleep(time.Millisecond)
	}
	return &game
}

// Effectue une requête au serveur, retourne une map du JSON retourné par le serveur.
func (c Game) Req(ob map[string]interface{}) (map[string]interface{}, error) {
	msg, err := json.Marshal(ob)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.Url, bytes.NewBuffer(msg))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var repMap map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&repMap); err != nil {
		return nil, err
	}
	return repMap, nil
}

// Convertit un array d'interfaces en array de strings
func IntToStr(intefaceArray []interface{}, strArray []string) {
	for i, v := range intefaceArray {
		switch v.(type) {
		case string:
			strArray[i] = v.(string)
		default:
			strArray[i] = ""
		}
	}
}
