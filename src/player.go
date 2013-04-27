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

import ()

type Player struct {
	*Sprite
}

// IncreaseVelocityTowardsEntity computes a new velocity vector for p.
func (p *Player) IncreaseVelocityTowardsEntity(s *Sprite) {
	var (
		pc = p.Centroid()
		sc = s.Centroid()
	)
	p.VelocityX = (sc.X - pc.X) / 10.0
	p.VelocityY = (sc.Y - pc.Y) / 10.0
}
