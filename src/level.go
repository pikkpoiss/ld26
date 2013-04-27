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
	"log"
	"math"
)

// Need to load a level, call something like this:
// twodee.LoadTiledMap(system, level, "examples/complex/levels/level01.json");

type SpatialClass int

const (
	CIRCLE SpatialClass = iota
	BOX
	DYNAMIC
)

// Level describes a particular set of static and dynamic sprites that make up a
// particular map.
type Level struct {
	System      *twodee.System
	Entities    map[SpatialClass][]SpatialVisible
	levelBounds twodee.Rectangle
	Player      *Player
}

// NewLevel constructs a Level struct and returns it.
func NewLevel(s *twodee.System) *Level {
	return &Level{
		System: s,
		Entities: map[SpatialClass][]SpatialVisible{
			CIRCLE:  make([]SpatialVisible, 0),
			BOX:     make([]SpatialVisible, 0),
			DYNAMIC: make([]SpatialVisible, 0),
		},
	}
}

// Create generates a new sprite for the indicated tileset and appends it to our
// set of level sprites. It gets called when the level loader recognizes a tile.
func (l *Level) Create(tileset string, index int, x, y, w, h float64) {
	// Need to check the tileset and index to determine what to do.
	// Should create sprites for most objects
	// Keep track of player sprite and just mark starting location.
	switch tileset {
	case "sprites32":
		var sprite = &Sprite{l.System.NewSprite(tileset, x, y, w, h, index)}
		sprite.SetFrame(index)
		if index == 0 {
			l.Player = &Player{sprite}
		} else {
			l.Entities[CIRCLE] = append(l.Entities[CIRCLE], sprite)
		}
	case "sprites16":
		var sprite = &Sprite{l.System.NewSprite(tileset, x, y, w, h, index)}
		sprite.SetFrame(index)
		l.Entities[BOX] = append(l.Entities[BOX], sprite)
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

// GetCollision returns whatever spatial s first collides with in the level.
func (l *Level) GetCollision(s twodee.Spatial) twodee.Spatial {
	r := s.Bounds()
	for _, eClass := range l.Entities {
		for _, e := range eClass {
			if s == e {
				continue
			}
			if r.Overlaps(e.Bounds()) {
				return e
			}
		}
	}
	return nil
}

// Draw iterates over all entities in the level and draws them.
func (l *Level) Draw() {
	for _, eClass := range l.Entities {
		for _, e := range eClass {
			e.Draw()
		}
	}
	l.Player.Draw()
}

// GetClosestEntity returns the closest CIRCLE type entity to the given entity.
func (l *Level) GetClosestEntity(s *Sprite) Spatial {
	p := s.Centroid()
	ld := math.MaxFloat64
	var ce SpatialVisible = nil
	for _, e := range l.Entities[CIRCLE] {
		if s == e {
			continue
		}
		ep := e.Centroid()
		if d := math.Hypot(p.X-ep.X, p.Y-ep.Y); d < ld {
			ld = d
			ce = e
		}
	}
	return ce
}
