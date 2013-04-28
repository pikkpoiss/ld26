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

type SpinOption int

const (
	CLOCKWISE SpinOption = iota
	COUNTER_CLOCKWISE
	NOSPIN
)

type Spatial interface {
	Centroid() twodee.Point
	twodee.Spatial
	twodee.Visible
}

type Moveable interface {
	Spatial
	MoveToward(s Spatial)
}

type Attachable interface {
	Spatial
	SetState(int)
	State() int
	Spin() SpinOption
}

type Colliding interface {
	Spatial
	Collision(p *Player)
}
