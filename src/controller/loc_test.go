package controller

import (
	"github.com/faiface/pixel/pixelgl"
	"os"
	"testing"
	"time"
)

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

func TestMain(m *testing.M) {
	pixelgl.Run(func() {
		os.Exit(m.Run())
	})
}

func TestNewResto(t *testing.T) {
	for _, v := range NewRestoTest {
		go func(l NewRestoTestStruct) {
			resto := NewResto(1280, 704, l.temps, l.accel, 1, l.pause, [][2]float64{{2.0, 3.0}},
				[]string{"e"}, []string{"p"}, []string{"d"})
			time.Sleep(501 * time.Millisecond)
			if l.tAttendu != resto.Temps {
				t.Errorf("Temps dans le restaurant après 0.5 seconde, accélération %v: attendu %v, reçu %v",
					l.accel, l.tAttendu, resto.Temps)
			}
		}(v)
	}
	time.Sleep(502 * time.Millisecond)
}
