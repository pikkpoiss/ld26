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
	BG_R int = 0
	BG_G int = 0
	BG_B int = 0
	BG_A int = 0
)

type Game struct {
	System *twodee.System
	Window *twodee.Window
	Camera *twodee.Camera
	exit   chan bool
}

func NewGame(sys *twodee.System, win *twodee.Window) (game *Game, err error) {
	game = &Game{
		System: sys,
		Window: win,
		Camera: twodee.NewCamera(0, 0, 71, 40),
		exit:   make(chan bool, 1),
	}
	if err = sys.Open(win); err != nil {
		err = fmt.Errorf("Couldn't open window: %v", err)
		return
	}
	game.handleKeys()
	game.handleClose()
	game.System.SetClearColor(BG_R, BG_G, BG_B, BG_A)
	return
}

func (g *Game) handleClose() {
	g.System.SetCloseCallback(func() int {
		g.exit <- true
		return 0
	})
}

func (g *Game) handleKeys() {
	g.System.SetKeyCallback(func(key int, state int) {
		switch {
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
}
