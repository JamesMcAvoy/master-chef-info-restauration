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
	pixelgl.Run(func() {
		os.Exit(m.Run())
	})
}

type NewRestoTestStruct struct {
	temps, accel int
	pause        bool
	tAttendu     int
}

var NewRestoTest = []NewRestoTestStruct{
	{0, 2, false, 1},
	{5, 1, true, 5},
	{0, 60, false, 30},
}

func TestNewResto(t *testing.T) {
	for _, v := range NewRestoTest {
		go func(l NewRestoTestStruct) {
			resto := NewResto(1280, 704, l.temps, l.accel, 1, l.pause, [][2]float64{{2.0, 3.0}},
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
