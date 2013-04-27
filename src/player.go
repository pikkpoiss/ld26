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
	"math"
)

type Player struct {
	*Sprite
}

// IncreaseVelocityTowardsEntity computes a new velocity vector for p.
func (p *Player) MoveToward(s *Sprite) {
	var (
		pc = p.Centroid()
		sc = s.Centroid()
		dx = (sc.X - pc.X)
		dy = (sc.Y - pc.Y)
		vx = math.Max(0.001, (100-math.Abs(dx))/1000.0)
		vy = math.Max(0.001, (100-math.Abs(dy))/1000.0)
	)
	if math.Signbit(dx) {
		p.VelocityX = -vx
	} else {
		p.VelocityX = vx
	}
	if math.Signbit(dy) {
		p.VelocityY = -vy
	} else {
		p.VelocityY = vy
	}
}
