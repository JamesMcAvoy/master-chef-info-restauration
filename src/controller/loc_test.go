package controller

import (
	"github.com/faiface/pixel/pixelgl"
	"os"
	"reflect"
	"testing"
	"time"
)

// Donne le thread principal à pixel
func TestMain(m *testing.M) {
	err := os.Chdir(os.Getenv("GOPATH") + "/src/github.com/JamesMcAvoy/resto")
	if err != nil {
		panic(err)
	}
	pixelgl.Run(func() {
		os.Exit(m.Run())
	})
}

type NewRestoTestStruct struct {
	temps, accel int
	tAttendu     int
}

var NewRestoTest = []NewRestoTestStruct{
	{0, 2, 1},
	{5, 1, 10},
	{0, 60, 30},
}

func TestNewResto(t *testing.T) {
	for _, v := range NewRestoTest {
		go func(l NewRestoTestStruct) {
			resto := NewResto(1280, 704, l.temps, l.accel, [][2]float64{{2.0, 3.0}},
				[]string{"e"}, []string{"p"}, []string{"d"}, []interface{}{})
			time.Sleep(501 * time.Millisecond)
			if l.tAttendu != resto.Temps {
				t.Errorf("Temps dans le restaurant après 0.5 seconde, accélération %v: attendu %v, reçu %v",
					l.accel, l.tAttendu, resto.Temps)
			}
		}(v)
	}
	time.Sleep(502 * time.Millisecond)
}

var EstOuvertTest = []struct {
	temps    int
	horaires [][2]float64
	expected bool
}{
	{100, nil, false},
	{100, [][2]float64{{0, 10}}, false},
	{100, [][2]float64{{90, 110}}, true},
	{100, [][2]float64{{90, 99}, {101, 130}}, false},
	{100, [][2]float64{{90, 95}, {99, 130}}, true},
}

func TestEstOuvert(t *testing.T) {
	resto := NewResto(1000, 1000, 0, 1, [][2]float64{}, []string{"e"}, []string{"p"}, []string{"d"}, []interface{}{})
	for _, v := range EstOuvertTest {
		resto.Temps = v.temps
		resto.Horaires = v.horaires
		if resto.EstOuvert() != v.expected {
			t.Errorf("resto.EstOuvert(): temps %v, horaires %v, attendu %t, reçu %t", v.temps, v.horaires, v.expected, resto.EstOuvert())
		}
	}
}

var RépartitTest = []struct {
	width, height, nb int
	output            [][4]int
}{
	{1000, 500, 0, nil},
	{1000, 500, 1, [][4]int{{0, 0, 1000, 500}}},
	{1000, 500, 2, [][4]int{{0, 0, 500, 500}, {500, 0, 1000, 500}}},
	{1000, 500, 3, [][4]int{{0, 0, 500, 250}, {0, 250, 500, 500}, {500, 0, 1000, 500}}},
	{1000, 1000, 4, [][4]int{{0, 0, 500, 500}, {0, 500, 500, 1000}, {500, 0, 1000, 500}, {500, 500, 1000, 1000}}},
	{1000, 1000, 10, [][4]int{
		{0, 0, 333, 250}, {0, 250, 333, 500}, {0, 500, 333, 750}, {0, 750, 333, 1000},
		{333, 0, 666, 333}, {333, 333, 666, 666}, {333, 666, 666, 999},
		{666, 0, 999, 333}, {666, 333, 999, 666}, {666, 666, 999, 999},
	}},
	{1000, 1000, 12, [][4]int{
		{0, 0, 250, 333}, {0, 333, 250, 666}, {0, 666, 250, 999},
		{250, 0, 500, 333}, {250, 333, 500, 666}, {250, 666, 500, 999},
		{500, 0, 750, 333}, {500, 333, 750, 666}, {500, 666, 750, 999},
		{750, 0, 1000, 333}, {750, 333, 1000, 666}, {750, 666, 1000, 999},
	}},
	{1000, 1000, 13, [][4]int{
		{0, 0, 250, 250}, {0, 250, 250, 500}, {0, 500, 250, 750}, {0, 750, 250, 1000},
		{250, 0, 500, 333}, {250, 333, 500, 666}, {250, 666, 500, 999},
		{500, 0, 750, 333}, {500, 333, 750, 666}, {500, 666, 750, 999},
		{750, 0, 1000, 333}, {750, 333, 1000, 666}, {750, 666, 1000, 999},
	}},
}

func TestRépartit(t *testing.T) {
	for _, v := range RépartitTest {
		output := Répartit(v.width, v.height, v.nb)
		if !reflect.DeepEqual(output, v.output) {
			t.Errorf("Répartit(%v, %v, %v):\nattendu:  %v\nreçu: %v", v.width, v.height, v.nb, v.output, output)
		}
	}
}
