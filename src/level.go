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

// Need to load a level, call something like this:
// twodee.LoadTiledMap(system, level, "examples/complex/levels/level01.json");

// Level describes a particular set of static and dynamic sprites that make up a
// particular map.
type Level struct {
	System      *twodee.System
	Static      []twodee.SpatialVisible
	Dynamic     []twodee.SpatialVisibleChanging
	levelBounds twodee.Rectangle
}

// Create generates a new sprite for the indicated tileset and appends it to our
// set of level sprites. It gets called when the level loader recognizes a tile.
func (l *Level) Create(tileset string, index int, x, y, w, h float64) {
	// Need to check the tileset and index to determine what to do.
	// Should create sprites for most objects
	// Keep track of player sprite and just mark starting location.
	switch tileset {
	case "tilegame":
		var sprite = l.System.NewSprite(tileset, x, y, w, h, index)
		sprite.SetFrame(index)
		l.Static = append(l.Static, sprite)
	case "character-textures":
		var sprite = l.System.NewSprite(tileset, x, y, w, h, index)
		var creature = NewCreature(sprite)
		creature.SetFrame(index)
		l.Dynamic = append(l.Dynamic, creature)
		l.player = creature
	default:
		log.Printf("Tileset: %v %v\n", tileset, index)
		log.Printf("Dim: %v %v %v %v\n", x, y, w, h)
	}
}

// SetBounds stores the size of this particular level.
func (l *Level) SetBounds(rect twodee.Rectangle) {
	l.levelBounds = rect
}

// GetBounds returns the size of this particular level.
func (l *Level) GetBounds() twodee.Rectangle {
	return l.levelBounds
}
