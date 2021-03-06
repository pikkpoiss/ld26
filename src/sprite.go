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
)

type Sprite struct {
	*twodee.Sprite
	state int
}

func NewSprite(sprite *twodee.Sprite) *Sprite {
	return &Sprite{
		Sprite: sprite,
		state: STATE_NORMAL,
	}
}

// Centroid returns the Point at the center of this sprite.
func (s *Sprite) Centroid() twodee.Point {
	var b = s.Sprite.Bounds()
	var p = twodee.Pt(b.Min.X+(b.Max.X-b.Min.X)/2.0, b.Min.Y+(b.Max.Y-b.Min.Y)/2.0)
	return p
}

func (s *Sprite) SetState(state int) {
	s.state = state
}

func (s *Sprite) State() int {
	return s.state
}

func (s *Sprite) Collision(p *Player) {
	// No-op in the default case
}
