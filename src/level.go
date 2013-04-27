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

type Level struct {
}

// Need to load a level, call something like this:
// twodee.LoadTiledMap(system, level, "examples/complex/levels/level01.json");

// Gets called when the level loader recognizes a tile
func (l *Level) Create(tileset string, index int, x, y, w, h float64) {
	// Need to check the tileset and index to determine what to do.
	// Should create sprites for most objects
	// Keep track of player sprite and just mark starting location.
}


