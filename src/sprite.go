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
}

// Centroid returns the Point at the center of this sprite.
func (s *Sprite) Centroid() twodee.Point {
	var b = s.Sprite.GlobalBounds()
	return twodee.Pt((b.Max.X-b.Min.X)/2, (b.Max.Y-b.Min.Y)/2)
}
