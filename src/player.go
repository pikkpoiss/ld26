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
	"time"
)

const (
	PLAYER_NORMAL = iota
	PLAYER_DAMAGED
)

type Player struct {
	*Sprite
	start   twodee.Point
	Elapsed time.Duration
	Damage  float64
	Won     bool
	frames  map[int]*twodee.Animation
	state   int
}

func NewPlayer(sprite *Sprite) *Player {
	return &Player{
		Sprite:  sprite,
		Elapsed: 0,
		Damage:  0,
		start:   twodee.Pt(sprite.X(), sprite.Y()),
		frames: map[int]*twodee.Animation{
			PLAYER_NORMAL: twodee.Anim([]int{0, 1, 2, 6}, 4),
			PLAYER_DAMAGED: twodee.Anim([]int{7, 8, 9, 10}, 4),
		},
		state: PLAYER_NORMAL,
	}
}

// IncreaseVelocityTowardsEntity computes a new velocity vector for p.
func (p *Player) MoveToward(s Spatial) {
	var (
		pc = p.Centroid()
		sc = s.Centroid()
		dx = (sc.X - pc.X)
		dy = (sc.Y - pc.Y)
		h  = math.Hypot(dx, dy)
		vx = math.Max(1, 5-h) * 0.2 * dx / h
		vy = math.Max(1, 5-h) * 0.2 * dy / h
	)
	p.VelocityX += (vx - p.VelocityX) / 40
	p.VelocityY += (vy - p.VelocityY) / 40
}

// GravitateToward computes a new velocity vector for p, taking into account a
// circulation effect.
func (p *Player) GravitateToward(s Attachable) {
	var (
		pc  = p.Centroid()
		sc  = s.Centroid()
		avx = sc.X - pc.X
		avy = sc.Y - pc.Y
		d   = math.Hypot(avx, avy)
	)
	// Normalize vector and include sensible constraints.
	avx = avx / d
	avy = avy / d
	av := twodee.Pt(math.Max(1, 5-d)*0.2*avx, math.Max(1, 5-d)*0.2*avy)

	// There are two possible orthogonal 'circulation' vectors.
	cv1 := twodee.Pt(-av.Y, av.X)
	cv2 := twodee.Pt(av.Y, -av.X)
	cv := cv1

	// Decide which vector to use based on the 'spin' of the spatial.
	switch s.Spin() {
	case CLOCKWISE:
		if av.X > 0 {
			// On the left side of the well, our Y component should be positive.
			cv = cv1
		} else {
			// On the right side of the well, our Y component should be negative.
			cv = cv2
		}
	case COUNTER_CLOCKWISE:
		if av.X > 0 {
			// On the left side of the well, our Y component should be negative.
			cv = cv2
		} else {
			// On the right side of the well, our Y component should be positive
			cv = cv1
		}
	default:
		// Compute whichever circulation vector is closer to our present vector.
		// cos(theta) = A -dot- B / ||A||*||B||
		dp1 := p.VelocityX*cv1.X + p.VelocityY*cv1.Y
		denom := math.Sqrt(p.VelocityX*p.VelocityX + p.VelocityY*p.VelocityY)
		theta1 := dp1 / denom
		dp2 := p.VelocityX*cv2.X + p.VelocityY*cv2.Y
		theta2 := dp2 / denom
		if theta1 >= theta2 {
			cv = cv1
		} else {
			cv = cv2
		}
	}

	// Now do some vector addition.
	fv := twodee.Pt(av.X+cv.X, av.Y+cv.Y)
	p.VelocityX += (fv.X - p.VelocityX) / 40
	p.VelocityY += (fv.Y - p.VelocityY) / 40
}

const (
	CT = iota
	CB
	CL
	CR
)

func (p *Player) Injure(amt float64) {
	p.Damage += amt
	p.state = PLAYER_DAMAGED
}

func (p *Player) Bounce(t Spatial) {
	bp := p.Bounds()
	bt := t.Bounds()
	dist := 1000.0
	coll := CT
	if d := math.Abs(bp.Max.X - bt.Min.X); d < dist {
		dist = d
		coll = CL
	}
	if d := math.Abs(bp.Min.X - bt.Max.X); d < dist {
		dist = d
		coll = CR
	}
	if d := math.Abs(bp.Max.Y - bt.Min.Y); d < dist {
		dist = d
		coll = CB
	}
	if d := math.Abs(bp.Min.Y - bt.Max.Y); d < dist {
		dist = d
		coll = CT
	}
	switch coll {
	case CL:
		p.VelocityX = -math.Abs(p.VelocityX)
		p.MoveTo(twodee.Pt(t.X()-p.Width(), p.Y()))
	case CR:
		p.VelocityX = math.Abs(p.VelocityX)
		p.MoveTo(twodee.Pt(t.X()+t.Width(), p.Y()))
	case CT:
		p.VelocityY = math.Abs(p.VelocityY)
		p.MoveTo(twodee.Pt(p.X(), t.Y()+t.Height()))
	case CB:
		p.VelocityY = -math.Abs(p.VelocityY)
		p.MoveTo(twodee.Pt(p.X(), t.Y()-p.Height()))
	}
}

func (p *Player) Reset() {
	p.VelocityX = 0
	p.VelocityY = 0
	p.MoveTo(p.start)
	p.Elapsed = 0
	p.Damage = 0
}

func (p *Player) Update() {
	p.Elapsed += time.Second / time.Duration(UPDATE_HZ)
	p.Sprite.SetFrame(p.frames[p.state].Next())
	p.Sprite.Update()
	p.state = PLAYER_NORMAL
}

func (p *Player) SetWon(w bool) {
	p.Won = w
}
