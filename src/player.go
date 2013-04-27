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

type Player struct {
	*Sprite
	start twodee.Point
}

func NewPlayer(sprite *Sprite) *Player {
	return &Player{
		Sprite: sprite,
		start:  twodee.Pt(sprite.X(), sprite.Y()),
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

// Ceases all movement; sets velocities to 0.
func (p *Player) SignalCollision() {
	p.VelocityX = -p.VelocityX
	p.VelocityY = -p.VelocityY
}

func (p *Player) Reset() {
	p.VelocityX = 0
	p.VelocityY = 0
	p.MoveTo(p.start)
}
