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

type Player struct {
	*Sprite
	start   twodee.Point
	Elapsed time.Duration
	Damage  float64
	Won bool
}

func NewPlayer(sprite *Sprite) *Player {
	return &Player{
		Sprite:  sprite,
		Elapsed: 0,
		Damage:  0,
		start:   twodee.Pt(sprite.X(), sprite.Y()),
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
func (p *Player) GravitateToward(s Spatial) {
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

	// Calculate an orthogonal counter-clockwise 'circulation' vector.
	cv := twodee.Pt(-av.Y, av.X)
	if av.X > 0 {
		// On the left side of the well, our Y component should be negative.
		if cv.Y > 0 {
			cv.X = -cv.X
			cv.Y = -cv.Y
		}
	} else {
		// On the right side of the well, our Y component should be positive
		if cv.Y < 0 {
			cv.X = -cv.X
			cv.Y = -cv.Y
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
}

func (p *Player) Update() {
	p.Elapsed += time.Second / time.Duration(UPDATE_HZ)
	p.Sprite.Update()
}

func (p *Player) SetWon(w bool) {
	p.Won = w
}
