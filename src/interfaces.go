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

type Spatial interface {
	Centroid() twodee.Point
	twodee.Spatial
}

type SpatialVisible interface {
	Spatial
	twodee.Visible
}

type SpatialVisibleMovable interface {
	Spatial
	twodee.Visible
	MoveToward(s Spatial)
}

type Stateful interface {
	SetState(int)
	State() int
}

type SpatialVisibleStateful interface {
	SpatialVisible
	Stateful
}

