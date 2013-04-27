// Copyright 2013 Arne Roomann-Kurrik + Wes Goodman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"../lib/twodee"
	"fmt"
	"log"
	"time"
)

const (
	UPDATE_HZ int = 60
	PAINT_HZ  int = 60
	BG_R      int = 0
	BG_G      int = 0
	BG_B      int = 0
	BG_A      int = 0
)

const (
	STATE_SPLASH = iota
	STATE_GAME
)

const (
	STATE_NORMAL = iota
	STATE_CLOSEST
	STATE_ATTACHED
)

type Game struct {
	System  *twodee.System
	Window  *twodee.Window
	Camera  *twodee.Camera
	Level   *Level
	Splash  *twodee.Sprite
	state   int
	exit    chan bool
	closest SpatialVisibleStateful
}

func NewGame(sys *twodee.System, win *twodee.Window) (game *Game, err error) {
	var (
		font *twodee.Font
	)
	game = &Game{
		System: sys,
		Window: win,
		Camera: twodee.NewCamera(0, 0, 71, 40),
		state:  STATE_SPLASH,
		exit:   make(chan bool, 1),
	}
	if err = sys.Open(win); err != nil {
		err = fmt.Errorf("Couldn't open window: %v", err)
		return
	}
	if err = sys.LoadTexture("splash", "data/splash.png", twodee.IntNearest, 1136); err != nil {
		err = fmt.Errorf("Couldn't load texture: %v", err)
		return
	}
	if font, err = twodee.LoadFont("data/slkscr.ttf", 24); err != nil {
		err = fmt.Errorf("Couldn't load font: %v", err)
		return
	}
	game.Splash = game.System.NewSprite("splash", 0, 0, 71, 40, 0)
	game.Splash.SetTextureHeight(640)
	game.handleKeys()
	game.handleClose()
	game.System.SetFont(font)
	game.System.SetClearColor(BG_R, BG_G, BG_B, BG_A)
	game.Level = NewLevel(game.System)
	twodee.LoadTiledMap(game.System, game.Level, "data/level-dev.json")
	return
}

func (g *Game) handleClose() {
	g.System.SetCloseCallback(func() int {
		g.exit <- true
		return 0
	})
}

func (g *Game) checkKeys() {
	switch {
	case g.System.Key(twodee.KeySpace) == 1:
		// Handle player shit
		switch g.state {
		case STATE_GAME:
			g.closest.SetState(STATE_ATTACHED)
			g.Level.Player.MoveToward(g.closest)
		}
	}
}

func (g *Game) handleKeys() {
	g.System.SetKeyCallback(func(key int, state int) {
		switch {
		case key == twodee.KeySpace && state == 0:
			switch g.state {
			case STATE_SPLASH:
				g.state = STATE_GAME
			}
		case state == 0:
			return
		case key == twodee.KeyEsc:
			g.exit <- true
		default:
			log.Printf("Key: %v, State: %v\n", key, state)
		}
	})
}

func (g *Game) Run() (err error) {
	go func() {
		update := time.NewTicker(time.Second / time.Duration(UPDATE_HZ))
		for true {
			<-update.C
			if g.closest != nil {
				g.closest.SetState(STATE_NORMAL)
			}
			g.closest = g.Level.GetClosestAttachable(g.Level.Player.Sprite)
			g.closest.SetState(STATE_CLOSEST)
			g.checkKeys()
			if g.Level.Player != nil {
				g.Level.Player.Update()
				if c := g.Level.GetCollision(g.Level.Player); c != nil {
					g.Level.Player.SignalCollision()
				}
			}
		}
	}()
	running := true
	paint := time.NewTicker(time.Second / time.Duration(PAINT_HZ))
	for running == true {
		<-paint.C
		g.System.Paint(g)
		select {
		case <-g.exit:
			paint.Stop()
			running = false
		default:
		}
	}
	return
}

func (g *Game) Draw() {
	g.Camera.SetProjection()
	if g.state == STATE_SPLASH {
		g.Splash.Draw()
	} else {
		g.Level.Draw()
	}
}
