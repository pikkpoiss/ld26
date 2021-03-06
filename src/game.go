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
	STATE_SUMMARY
	STATE_WIN
	STATE_DIED
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
	Scores       []*Score
	CurrentLevel int
	Splash       *twodee.Sprite
	SelectMenu   *LevelSelect
	OptionMenu   *Menu
	Summary      *Summary
	state        int
	laststate    int
	lastpaint    time.Time
	exit         chan bool
	closest      Attachable
}

func NewGame(sys *twodee.System, win *twodee.Window) (game *Game, err error) {
	game = &Game{
		System: sys,
		Window: win,
		Camera: twodee.NewCamera(0, 0, 71, 40),
		Levels: []string{
			"data/level-01.json",
			"data/level-02.json",
			"data/level-03.json",
			"data/level-04.json",
			"data/level-05.json",
			"data/level-06.json",
		},
		Scores:       make([]*Score, 6),
		CurrentLevel: 0,
		state:        STATE_SPLASH,
		exit:         make(chan bool, 1),
	}
	if err = sys.Open(win); err != nil {
		err = fmt.Errorf("Couldn't open window: %v", err)
		return
	}
	if err = sys.LoadTexture("splash", "data/splash.png", twodee.IntNearest, 1136); err != nil {
		err = fmt.Errorf("Couldn't load texture: %v", err)
		return
	}
	if game.Font, err = twodee.LoadFont("data/slkscr.ttf", 32); err != nil {
		err = fmt.Errorf("Couldn't load font: %v", err)
		return
	}
	game.Splash = game.System.NewSprite("splash", 0, 0, 71, 40, 0)
	game.Splash.SetTextureHeight(640)
	game.Splash.SetFrame(0)
	game.handleKeys()
	game.handleClose()
	game.System.SetClearColor(BG_R, BG_G, BG_B, BG_A)
	game.SelectMenu = NewLevelSelect(game.System, len(game.Levels))
	if twodee.LoadTiledMap(game.System, game.SelectMenu, "data/menu-select.json"); err != nil {
		err = fmt.Errorf("Couldn't load select menu: %v", err)
		return
	}
	game.SelectMenu.SetSelectable(0, true)
	game.OptionMenu = NewMenu(game.System)
	log.Printf("Loading option menu\n")
	if err = twodee.LoadTiledMap(game.System, game.OptionMenu, "data/menu-options.json"); err != nil {
		err = fmt.Errorf("Couldn't load options menu: %v", err)
		return
	}
	game.OptionMenu.SetSelectable(0, true)
	game.OptionMenu.SetSelectable(1, false)

	game.Summary = NewSummary(game.System, game.Font, game.Window)
	if twodee.LoadTiledMap(game.System, game.Summary, "data/menu-summary.json"); err != nil {
		err = fmt.Errorf("Couldn't load summary menu %v", err)
		return
	}

	Play()
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

func (g *Game) handleMenuKey(menu Navigatable, key int, state int) {
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
			case key == twodee.KeyEnter && state == 1:
				fallthrough
			case key == twodee.KeySpace && state == 1:
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
			case key == twodee.KeyEsc && state == 1:
				g.state = g.laststate
			}
		case STATE_SELECT:
			g.handleMenuKey(g.SelectMenu, key, state)
			switch {
			case key == twodee.KeyEnter && state == 1:
				fallthrough
			case key == twodee.KeySpace && state == 1:
				g.SetLevel(g.SelectMenu.GetSelection())
			case key == twodee.KeyEsc && state == 1:
				g.OptionMenu.SetSelection(0)
				g.OptionMenu.SetSelectable(1, false)
				g.laststate = g.state
				g.state = STATE_OPTION
			}
		case STATE_DIED:
			switch {
			case state == 1:
				g.state = STATE_GAME
			}
		case STATE_WIN:
			fallthrough
		case STATE_SPLASH:
			switch {
			case key == twodee.KeyEnter && state == 1:
				fallthrough
			case key == twodee.KeySpace && state == 1:
				g.state = STATE_SELECT
			case key == twodee.KeyEsc && state == 1:
				g.OptionMenu.SetSelection(0)
				g.OptionMenu.SetSelectable(1, false)
				g.laststate = g.state
				g.state = STATE_OPTION
			}
		case STATE_GAME:
			switch {
			case key == twodee.KeyEsc && state == 1:
				g.OptionMenu.SetSelectable(1, true)
				g.OptionMenu.SetSelection(1)
				g.laststate = g.state
				g.state = STATE_OPTION
			case key == 80 && state == 1: // p
				//g.handleWin()
			default:
			}
		case STATE_SUMMARY:
			switch {
			case state == 1:
				if g.CurrentLevel+1 == len(g.Levels) {
					g.state = STATE_WIN
					g.Splash.SetFrame(2)
				} else {
					g.state = STATE_SELECT
				}
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
	g.SelectMenu.SetSelectable(g.CurrentLevel, true)
	g.SelectMenu.SetSelection(g.CurrentLevel)
	g.state = STATE_GAME
}

func (g *Game) handleWin() {
	score := &Score{
		Time:   g.Level.Player.Elapsed,
		Damage: g.Level.Player.Damage,
		Level:  g.CurrentLevel + 1,
	}
	g.Level.GetStars(score)
	g.Summary.SetMetrics(score)
	if score.BetterThan(g.Scores[g.CurrentLevel]) {
		g.Scores[g.CurrentLevel] = score
		g.SelectMenu.SetScores(g.Scores)
	}
	g.state = STATE_SUMMARY
	g.SelectMenu.SetSelectable(g.CurrentLevel+1, true)
	g.SelectMenu.NextSelection()
}

func (g *Game) Run() (err error) {
	go func() {
		update := time.NewTicker(time.Second / time.Duration(UPDATE_HZ))
		for true {
			<-update.C
			if g.Level == nil {
				continue
			}
			g.Level.Update()
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
				if g.Level.Player.Won {
					g.handleWin()
				}
			}
			b := g.Level.Player.Bounds()
			c := g.Camera.Bounds()
			if b.Max.X < c.Min.X || b.Min.X > c.Max.X ||
				b.Max.Y < c.Min.Y || b.Min.Y > c.Max.Y {
				// Player is offscreen so damage them
				g.Level.Player.Injure(0.005)
			}
			if g.Level.Player.Damage > 1 {
				g.state = STATE_DIED
				g.Splash.SetFrame(3)
				g.Level.Restart()
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
		now    = time.Now()
		//fps    = 1.0 / now.Sub(g.lastpaint).Seconds()
		option = false
		state  = g.state
	)
	g.Camera.SetProjection()
	if state == STATE_OPTION {
		state = g.laststate
		option = true
	}
	switch state {
	case STATE_DIED:
		fallthrough
	case STATE_WIN:
		fallthrough
	case STATE_SPLASH:
		g.Splash.Draw()
	case STATE_SELECT:
		g.SelectMenu.Draw()
	case STATE_SUMMARY:
		g.Summary.Draw()
	case STATE_GAME:
		if g.Level != nil {
			g.Level.Draw()
			if g.Level.Player.Damage > 0 {
				g.Font.Printf(10, 40, "Damage %.2f", g.Level.Player.Damage)
			}
			g.Font.Printf(10, 10, "%.1f seconds", g.Level.Player.Elapsed.Seconds())
			//g.Font.Printf(0, 10, "FPS %.1f", fps)
		}
	}
	if option {
		if g.state == STATE_OPTION {
			g.OptionMenu.Draw()
		}
	}
	g.lastpaint = now
}
