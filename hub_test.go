package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractRoom(t *testing.T) {
	cases := []struct {
		rooms []Roomer
		in    string
		want  Roomer
	}{
		{[]Roomer{}, "hoge", &Room{name: "hoge"}},
		{[]Roomer{&Room{name: "hoge"}}, "hoge", &Room{name: "hoge"}},
		{[]Roomer{&Room{name: "hoge"}}, "new", &Room{name: "new"}},
	}

	// set up
	root, err := os.Getwd()
	if err != nil {
		t.Errorf("failed pwd: %v", err)
	}
	for _, c := range cases {
		hub := Hub{
			Rooms:       c.rooms,
			historyRoot: root,
		}
		// execute
		r, err := hub.ExtractRoom(c.in)
		if err != nil {
			t.Errorf("failed extract: %v", err)
		}
		if c.want.Name() != r.Name() {
			t.Errorf("not expected room\n e:%s \n a:%s", c.want.Name(), r.Name())
		}

		// tear down
		_, err = os.Stat(filepath.Join(hub.historyRoot, r.Name()))
		if os.IsNotExist(err) {
			continue
		}

		err = os.Remove(filepath.Join(hub.historyRoot, r.Name()))
		if err != nil {
			t.Errorf("failed rm history file: %v", err)
		}
	}
}

func TestClose(t *testing.T) {
	cases := []struct {
		rooms []Roomer
		in    *Room
		want  int
		isErr bool
	}{
		{[]Roomer{&Room{name: "hoge"}},
			&Room{name: "hoge"},
			0,
			false,
		},
		{[]Roomer{},
			&Room{name: "hoge"},
			0,
			true,
		},
		{[]Roomer{&Room{name: "hoge"}},
			&Room{name: "fuga"},
			1,
			true,
		},
		{[]Roomer{&Room{name: "hoge"}, &Room{name: "hoge2"}},
			&Room{name: "hoge"},
			1,
			false,
		},
	}

	for _, c := range cases {
		hub := Hub{
			Rooms: c.rooms,
		}
		err := hub.Close(c.in)

		if len(hub.Rooms) != c.want {
			t.Errorf("not expected room num\n e:%d\n a:%d", c.want, len(hub.Rooms))
		}
		if c.isErr {
			if err == nil {
				t.Errorf("not occur expected error")
			}
		}
	}
}
