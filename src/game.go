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
	STATE_SELECT
	STATE_OPTION
	STATE_GAME
)

const (
	STATE_NORMAL = iota
	STATE_CLOSEST
	STATE_ATTACHED
)

type Game struct {
	System       *twodee.System
	Window       *twodee.Window
	Camera       *twodee.Camera
	Font         *twodee.Font
	Level        *Level
	Levels       []string
	CurrentLevel int
	Splash       *twodee.Sprite
	SelectMenu   *Menu
	OptionMenu   *Menu
	state        int
	laststate    int
	lastpaint    time.Time
	exit         chan bool
	closest      Stateful
}

func NewGame(sys *twodee.System, win *twodee.Window) (game *Game, err error) {
	game = &Game{
		System: sys,
		Window: win,
		Camera: twodee.NewCamera(0, 0, 71, 40),
		Levels: []string{
			"data/level-dev.json",
		},
		state: STATE_SPLASH,
		exit:  make(chan bool, 1),
	}
	if err = sys.Open(win); err != nil {
		err = fmt.Errorf("Couldn't open window: %v", err)
		return
	}
	if err = sys.LoadTexture("splash", "data/splash.png", twodee.IntNearest, 1136); err != nil {
		err = fmt.Errorf("Couldn't load texture: %v", err)
		return
	}
	if game.Font, err = twodee.LoadFont("data/slkscr.ttf", 24); err != nil {
		err = fmt.Errorf("Couldn't load font: %v", err)
		return
	}
	game.Splash = game.System.NewSprite("splash", 0, 0, 71, 40, 0)
	game.Splash.SetTextureHeight(640)
	game.handleKeys()
	game.handleClose()
	game.System.SetClearColor(BG_R, BG_G, BG_B, BG_A)
	game.SelectMenu = NewMenu(game.System)
	if twodee.LoadTiledMap(game.System, game.SelectMenu, "data/menu-select.json"); err != nil {
		err = fmt.Errorf("Couldn't load select menu: %v", err)
		return
	}
	game.SelectMenu.SetSelectable(0, true)
	game.SelectMenu.SetSelectable(1, true)
	game.OptionMenu = NewMenu(game.System)
	log.Printf("Loading option menu\n")
	if err = twodee.LoadTiledMap(game.System, game.OptionMenu, "data/menu-options.json"); err != nil {
		err = fmt.Errorf("Couldn't load options menu: %v", err)
		return
	}
	game.OptionMenu.SetSelectable(0, true)
	game.OptionMenu.SetSelectable(1, false)
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
			// g.Level.Player.MoveToward(g.closest)
			g.Level.Player.GravitateToward(g.closest)
		}
	}
}

func (g *Game) handleMenuKey(menu *Menu, key int, state int) {
	switch {
	case state == 0:
		return
	case key == twodee.KeyLeft:
		fallthrough
	case key == twodee.KeyUp:
		menu.PrevSelection()
	case key == twodee.KeyRight:
		fallthrough
	case key == twodee.KeyDown:
		menu.NextSelection()
	}
}

func (g *Game) handleKeys() {
	g.System.SetKeyCallback(func(key int, state int) {
		switch g.state {
		case STATE_OPTION:
			g.handleMenuKey(g.OptionMenu, key, state)
			switch {
			case key == twodee.KeyEnter && state == 0:
				fallthrough
			case key == twodee.KeySpace && state == 0:
				sel := g.OptionMenu.GetSelection()
				switch sel {
				case 0:
					g.exit <- true
				case 1:
					if g.Level != nil {
						g.Level.Restart()
					}
					g.state = g.laststate
				}
			case key == twodee.KeyEsc && state == 0:
				g.state = g.laststate
			}
		case STATE_SELECT:
			g.handleMenuKey(g.SelectMenu, key, state)
			switch {
			case key == twodee.KeyEnter && state == 0:
				fallthrough
			case key == twodee.KeySpace && state == 0:
				g.SetLevel(g.SelectMenu.GetSelection())
			case key == twodee.KeyEsc && state == 0:
				g.OptionMenu.SetSelectable(1, false)
				g.laststate = g.state
				g.state = STATE_OPTION
			}
		case STATE_SPLASH:
			switch {
			case key == twodee.KeyEnter && state == 0:
				fallthrough
			case key == twodee.KeySpace && state == 0:
				g.state = STATE_SELECT
			case key == twodee.KeyEsc && state == 0:
				g.OptionMenu.SetSelectable(1, false)
				g.laststate = g.state
				g.state = STATE_OPTION
			}
		case STATE_GAME:
			switch {
			case key == twodee.KeyEsc && state == 0:
				g.OptionMenu.SetSelectable(1, true)
				g.OptionMenu.SetSelection(1)
				g.laststate = g.state
				g.state = STATE_OPTION
			}
		}
	})
}

func (g *Game) SetLevel(i int) {
	var (
		index = (i + len(g.Levels)) % len(g.Levels)
		path  = g.Levels[index]
		level = NewLevel(g.System)
	)
	twodee.LoadTiledMap(g.System, level, path)
	g.Level = level
	g.CurrentLevel = index
	g.state = STATE_GAME
}

func (g *Game) Run() (err error) {
	go func() {
		update := time.NewTicker(time.Second / time.Duration(UPDATE_HZ))
		for true {
			<-update.C
			if g.Level == nil {
				continue
			}
			if g.state != STATE_GAME {
				continue
			}
			if g.closest != nil {
				g.closest.SetState(STATE_NORMAL)
			}
			g.closest = g.Level.GetClosestAttachable(g.Level.Player.Sprite)
			g.closest.SetState(STATE_CLOSEST)
			g.checkKeys()
			if g.Level.Player != nil {
				g.Level.Player.Update()
				if c := g.Level.GetCollision(g.Level.Player); c != nil {
					c.Collision(g.Level.Player)
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
	var (
		now time.Time
		fps float64
	)
	now = time.Now()
	fps = 1.0 / now.Sub(g.lastpaint).Seconds()

	g.Camera.SetProjection()
	switch g.state {
	case STATE_SPLASH:
		g.Splash.Draw()
	case STATE_SELECT:
		g.SelectMenu.Draw()
	default:
		if g.Level != nil {
			g.Level.Draw()
		}
		g.Font.Printf(0, 10, "FPS %6.1f", fps)
		g.Font.Printf(0, 40, "%6.1f seconds", g.Level.Player.Elapsed.Seconds())
		if g.state == STATE_OPTION {
			g.OptionMenu.Draw()
		}
	}
	g.lastpaint = now
}
