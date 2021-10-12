package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"github.com/Shopify/go-lua"
)

//var color int = 0x00000000
var surface *sdl.Surface
var font *ttf.Font

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
	l.Register("print", func(l *lua.State) int {
		str := lua.CheckString(l, 1)
		x := lua.CheckInteger(l, 2)
		y := lua.CheckInteger(l, 3)
		text, err := font.RenderUTF8Blended(str, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err != nil {
			return 0
		}
		defer text.Free()
		text.Blit(nil, surface, &sdl.Rect{X: int32(x), Y: int32(y), W: 0, H: 0})
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

	if err = ttf.Init(); err != nil {
		return
	}
	defer ttf.Quit()

	if font, err = ttf.OpenFont("FiraSans-Regular.ttf", 32); err != nil {
		return
	}
	defer font.Close()

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

	text, err := font.RenderUTF8Blended(url, sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err != nil {
		return
	}
	defer text.Free()

	running := true
	for running {
		surface.FillRect(nil, 0)

		rect := sdl.Rect{0, 0, 800, 40}
		surface.FillRect(&rect, 0x00888888)
		text.Blit(nil, surface, &sdl.Rect{X: 10, Y: 0, W: 0, H: 0})

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

		sdl.Delay(16)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED {
					l.Global("mouse_pressed")
					if (l.IsFunction(-1)) {
						l.PushInteger(int(t.X))
						l.PushInteger(int(t.Y))
						err := l.ProtectedCall(2, 0, 0)
						if err != nil {
							log.Fatal(err)
						}
					}
				} else {
					l.Global("mouse_released")
					if (l.IsFunction(-1)) {
						l.PushInteger(int(t.X))
						l.PushInteger(int(t.Y))
						err := l.ProtectedCall(2, 0, 0)
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			case *sdl.QuitEvent:
				running = false
			}
		}
	}
}
