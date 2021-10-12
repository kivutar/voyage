package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/Shopify/go-lua"
)

//var color int = 0x00000000
var surface *sdl.Surface

func registerFuncs(l *lua.State) {
	l.Register("setColor", func(l *lua.State) int {
		return 0
	})
	l.Register("rectangle", func(l *lua.State) int {
		lua.CheckString(l, 1)
		x := lua.CheckInteger(l, 2)
		y := lua.CheckInteger(l, 3)
		w := lua.CheckInteger(l, 4)
		h := lua.CheckInteger(l, 5)
		rect := sdl.Rect{int32(x), int32(y), int32(w), int32(h)}
		surface.FillRect(&rect, 0x0000ffff)
		return 0
	})
}

func main() {
	url := os.Args[1]

	client := &gemini.Client{}
	ctx := context.Background()
	resp, err := client.Get(ctx, url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	program := string(bytes)

	l := lua.NewState()
	lua.OpenLibraries(l)
	registerFuncs(l)

	lua.DoString(l, program)

	l.Global("load")
	if (l.IsFunction(-1)) {
		err := l.ProtectedCall(0, 0, 0)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Voyage",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	surface, err = window.GetSurface()
	if err != nil {
		log.Fatal(err)
	}

	running := true
	for running {
		surface.FillRect(nil, 0)

		rect := sdl.Rect{0, 0, 800, 40}
		surface.FillRect(&rect, 0x00eeeeee)

		l.Global("update")
		if (l.IsFunction(-1)) {
			err := l.ProtectedCall(0, 0, 0)
			if err != nil {
				log.Fatal(err)
			}
		}

		l.Global("draw")
		if (l.IsFunction(-1)) {
			err := l.ProtectedCall(0, 0, 0)
			if err != nil {
				log.Fatal(err)
			}
		}

		window.UpdateSurface()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
