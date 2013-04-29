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
	"math"
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

func (g *GravityWell) Spin() SpinOption {
	return NOSPIN
}

const (
	ZONE_NORMAL = iota
)

type Zone struct {
	*Sprite
	frames map[int]*twodee.Animation
	state  int
}

func (z *Zone) Touches(p *Player) {
}

type HurtZone Zone

func NewHurtZone(sprite *twodee.Sprite) *HurtZone {
	return &HurtZone{
		Sprite: NewSprite(sprite),
	}
}

func (z *HurtZone) Collision(p *Player) {
	p.Injure(0.01)
}

type BounceZone Zone

func NewBounceZone(sprite *twodee.Sprite) *BounceZone {
	return &BounceZone{
		Sprite: NewSprite(sprite),
	}
}

func (z *BounceZone) Collision(p *Player) {
	var surface twodee.Point
	var cx, cy float64
	var cent = z.Centroid()
	var b = p.Bounds()
	switch z.Frame() {
	case 7: // Top left
		if b.Max.X < cent.X {
			if b.Min.Y > cent.Y {
				surface = twodee.Pt(1, 1)
			} else {
				surface = twodee.Pt(0, 1)
			}
		} else {
			surface = twodee.Pt(1, 0)
		}
		cx = z.X() - p.Width()
		cy = z.Y() + z.Height()
	case 8: // Bottom right
		if b.Min.X > cent.X {
			if b.Max.Y < cent.Y {
				surface = twodee.Pt(1, 1)
			} else {
				surface = twodee.Pt(0, 1)
			}
		} else {
			surface = twodee.Pt(1, 0)
		}
		cx = z.X() + z.Width()
		cy = z.Y() - p.Height()
	case 9: // Bottom left
		if b.Max.X < cent.X {
			if b.Max.Y < cent.Y {
				surface = twodee.Pt(-1, 1)
			} else {
				surface = twodee.Pt(0, 1)
			}
		} else {
			surface = twodee.Pt(1, 0)
		}
		cx = z.X() - p.Width()
		cy = z.Y() - p.Height()
	case 10: // Top right
		if b.Min.X > cent.X {
			if b.Min.Y > cent.Y {
				surface = twodee.Pt(-1, 1)
			} else {
				surface = twodee.Pt(0, 1)
			}
		} else {
			surface = twodee.Pt(1, 0)
		}
		cx = z.X() + z.Width()
		cy = z.Y() + z.Height()
	case 11: // Left
		surface = twodee.Pt(0, 1)
		cx = z.X() - p.Width()
		cy = 1000.0
	case 12: // Right
		surface = twodee.Pt(0, 1)
		cx = z.X() + z.Width()
		cy = 1000.0
	case 13: // Bottom
		surface = twodee.Pt(1, 0)
		cx = 1000.0
		cy = z.Y() - p.Height()
	case 14: // Top
		surface = twodee.Pt(1, 0)
		cx = 1000.0
		cy = z.Y() + z.Height()
	default:
		p.Bounce(z)
		return
	}
	if math.Abs(p.X()-cx) < math.Abs(p.Y()-cy) {
		p.MoveTo(twodee.Pt(cx, p.Y()))
	} else {
		p.MoveTo(twodee.Pt(p.X(), cy))
	}
	p.SetVelocity(z.reflect(p.Velocity(), surface))
}

func (z *BounceZone) dot(a twodee.Point, b twodee.Point) float64 {
	return a.X*b.X + a.Y*b.Y
}

func (z *BounceZone) reflect(v twodee.Point, l twodee.Point) twodee.Point {
	var (
		vl = z.dot(v, l)
		ll = z.dot(l, l)
		c  = 2 * vl / ll
	)
	return twodee.Pt(c*l.X-v.X, c*l.Y-v.Y)
}

type VictoryZone Zone

func NewVictoryZone(sprite *twodee.Sprite) *VictoryZone {
	return &VictoryZone{
		Sprite: NewSprite(sprite),
		frames: map[int]*twodee.Animation{
			ZONE_NORMAL: twodee.Anim([]int{0, 4, 5, 6}, 16),
		},
	}
}

func (z *VictoryZone) Collision(p *Player) {
	p.SetWon(true)
}

func (z *VictoryZone) Update() {
	z.SetFrame(z.frames[z.state].Next())
	z.Sprite.Update()
}
