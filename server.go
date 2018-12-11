package main

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Simple serveur temporaire en attendant de lier les 2 parties de l'application
func Serv(port int, accel int) {
	temps := 0
	go incTime(&temps, accel)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, temps)
	})
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func handle(w http.ResponseWriter, r *http.Request, temps int) {
	w.Header().Set("Content-Type", "application/json")
	var req map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Write([]byte("{\"type\": \"reponse\", \"error\": \"Impossible de parser le JSON en entrée\"}"))
		return
	}

	var rep map[string]interface{}
	switch req["type"] {
	case "bonjour":
		repBytes, err := ioutil.ReadFile("ressources/serveur/bonjour.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(repBytes, &rep)
	}
	rep["temps"] = temps
	repBytes, err := json.Marshal(rep)
	w.Write(repBytes)
}

func incTime(t *int, accel int) {
	tick := time.Tick(time.Second / time.Duration(accel))
	for {
		<-tick
		*t++
	}
}
