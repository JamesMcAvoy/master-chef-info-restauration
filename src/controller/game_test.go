package controller

import (
	"gopkg.in/h2non/gock.v1"
	"reflect"
	"testing"
)

var NewGameTest = []struct {
	width, height int
	url           string
	game          bool // Game a au moins 1 resto?
	apanic        bool // NewGame déclenche une panique?
}{
	{1280, 720, "http://url.net", true, false},
	{1280, 720, "http://url_qui_marche.pas", false, true},
}

func TestNewGame(t *testing.T) {
	defer gock.Off()
	gock.New("http://url.net").Get("/").Reply(200).File("../../testdata/bonjour.json")
	for _, v := range NewGameTest {
		defer func() {
			if r := recover(); (r == nil) != v.apanic {
				t.Errorf("NewGame(%q, %q, %q): a paniqué: %t attendu, %t reçu",
					v.width, v.height, v.url, v.apanic, (r == nil))
			}
		}()
		game := NewGame(v.width, v.height, v.url)
		if (len(game.Restos) > 0) != v.game {
			t.Errorf("NewGame(%q, %q, %q): au moins 1 restaurant: %t attendu, %t reçu",
				v.width, v.height, v.url, len(game.Restos) > 0, v.game)
		}
	}
}

var RequestTest = []struct {
	input  map[string]interface{}
	output map[string]interface{}
	err    error
}{
	{map[string]interface{}{"type": "bonjour"}, map[string]interface{}{"rep": "ok"}, nil},
	{
		map[string]interface{}{"type": map[string]interface{}{"oui": []interface{}{true, 43.4}}},
		map[string]interface{}{"rep": "ok"}, nil,
	},
}

func TestReq(t *testing.T) {
	defer gock.Off()
	gock.New("http://url.net").Get("/").Reply(200).JSON(map[string]string{"rep": "ok"})
	g := Game{
		Url:    "http://url.net",
		Restos: []*Resto{},
	}
	for _, v := range RequestTest {
		output, err := g.Req(v.input)
		if reflect.DeepEqual(output, v.output) && err != v.err {
			t.Errorf("Game.Req(%q) == %q, %q, veut %q, %q", v.input, output, err, v.output, v.err)
		}
	}
}

var IntToStrTest = []struct {
	input    []interface{}
	output   []string
	expected []string
}{
	{[]interface{}{"hop", "le", "test"}, make([]string, 3), []string{"hop", "le", "test"}},
	{[]interface{}{"hop", 2, "test"}, make([]string, 3), []string{"hop", "", "test"}},
}

func TestIntToStr(t *testing.T) {
	for _, v := range IntToStrTest {
		IntToStr(v.input, v.output)
		if !reflect.DeepEqual(v.output, v.expected) {
			t.Errorf("IntToStr(%s, []): %s attendu, %s reçu", v.input, v.expected, v.output)
		}
	}
}
