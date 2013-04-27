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
	"log"
)

// Mob represents a mobile sprite in the game.
type Mob struct {
	*Sprite
}

// NewMob creates a new Mob struct, setting its sprite field properly.
func NewMob(sprite *Sprite) *Mob {
	var m = &Mob{
		Sprite: sprite,
	}
	return m
}

type GravityWell struct {
	*Sprite
}

func NewGravityWell(sprite *twodee.Sprite) *GravityWell {
	return &GravityWell{
		Sprite: NewSprite(sprite),
	}
}

func (g *GravityWell) SetState(state int) {
	g.Sprite.SetState(state)
	switch state {
	case STATE_NORMAL:
		g.Sprite.SetFrame(3)
	case STATE_CLOSEST:
		g.Sprite.SetFrame(4)
	case STATE_ATTACHED:
		g.Sprite.SetFrame(5)
	}
}

type Zone struct {
	*Sprite
}

type VictoryZone Zone

func NewVictoryZone(sprite *twodee.Sprite) *VictoryZone {
	return &VictoryZone{
		Sprite: NewSprite(sprite),
	}
}

func (z *VictoryZone) Collision(p *Player) {
	log.Printf("Woop won")
}

