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
	"strconv"
)

// Need to load a level, call something like this:
// twodee.LoadTiledMap(system, level, "examples/complex/levels/level01.json");

type SpatialClass int

// Level describes a particular set of static and dynamic sprites that make up a
// particular map.
type Level struct {
	System      *twodee.System
	levelBounds twodee.Rectangle
	Player      *Player
	wells       []*GravityWell
	zones       []Colliding
	colliding   []Colliding
	startimes   []int
	stars       []*Sprite
}

// NewLevel constructs a Level struct and returns it.
func NewLevel(s *twodee.System) *Level {
	return &Level{
		System:    s,
		wells:     make([]*GravityWell, 0),
		startimes: []int{90, 60, 30},
	}
}

// Create generates a new sprite for the indicated tileset and appends it to our
// set of level sprites. It gets called when the level loader recognizes a tile.
func (l *Level) Create(tileset string, index int, x, y, w, h float64) {
	// Need to check the tileset and index to determine what to do.
	// Should create sprites for most objects
	// Keep track of player sprite and just mark starting location.
	switch tileset {
	case "stars":
		var sprite = NewSprite(l.System.NewSprite(tileset, x, y, w, h, index))
		sprite.SetFrame(index)
		l.stars = append(l.stars, sprite)
	case "sprites32":
		if index == 0 {
			var sprite = NewSprite(l.System.NewSprite(tileset, x, y, w, h, index))
			sprite.SetFrame(index)
			l.Player = NewPlayer(sprite)
		} else {
			var well = NewGravityWell(l.System.NewSprite(tileset, x, y, w, h, index))
			l.wells = append(l.wells, well)
			well.SetFrame(index)
		}
	case "sprites16":
		var sprite = l.System.NewSprite(tileset, x, y, w, h, index)
		sprite.SetFrame(index)
		var zone Colliding
		switch index {
		case 0:
			zone = NewVictoryZone(sprite)
			l.zones = append(l.zones, zone)
		case 2:
			zone = NewHurtZone(sprite)
			l.zones = append(l.zones, zone)
		case 15:
			zone = NewFlashingHurtZone(sprite)
			l.zones = append(l.zones, zone)
		default:
			zone = NewBounceZone(sprite)
			l.colliding = append(l.colliding, zone)
		}
	default:
		log.Printf("Tileset: %v %v\n", tileset, index)
		log.Printf("Dim: %v %v %v %v\n", x, y, w, h)
	}
}

// SetBounds stores the size of this particular level.
func (l *Level) Loaded(rect twodee.Rectangle, props map[string]string) {
	log.Printf("Loaded: %v\n", props)
	l.levelBounds = rect
	for k, v := range props {
		switch k {
		case "time1":
			l.startimes[0], _ = strconv.Atoi(v)
		case "time2":
			l.startimes[1], _ = strconv.Atoi(v)
		case "time3":
			l.startimes[2], _ = strconv.Atoi(v)
		}
	}
	log.Printf("Star times: %v\n", l.startimes)
}

// GetBounds returns the size of this particular level.
func (l *Level) GetBounds() twodee.Rectangle {
	return l.levelBounds
}

// GetCollision returns the closest colliding spatial in the level.
func (l *Level) GetCollision(s *Player) (out Colliding) {
	r := s.Bounds()
	ld := 1000.0
	sp := twodee.Pt(r.Max.X-r.Min.X/2.0, r.Max.Y-r.Min.Y/2.0)
	for _, e := range l.zones {
		if r.Overlaps(e.Bounds()) {
			e.Collision(s)
		}
	}
	for _, e := range l.colliding {
		if r.Overlaps(e.Bounds()) {
			ep := e.Centroid()
			if d := math.Hypot(ep.X-sp.X, ep.Y-sp.Y); d < ld {
				out = e
				ld = d
			}
		}
	}
	return
}

// Draw iterates over all entities in the level and draws them.
func (l *Level) Draw() {
	for _, e := range l.stars {
		e.Draw()
	}
	for _, e := range l.zones {
		e.Draw()
	}
	for _, e := range l.colliding {
		e.Draw()
	}
	for _, e := range l.wells {
		e.Draw()
	}
	l.Player.Draw()
}

// Resets the level
func (l *Level) Restart() {
	l.Player.Reset()
}

// GetClosestEntity returns the closest CIRCLE type entity to the given entity.
func (l *Level) GetClosestAttachable(s Spatial) Attachable {
	p := s.Centroid()
	ld := math.MaxFloat64
	var ce Attachable = nil
	for _, e := range l.wells {
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

func (l *Level) Update() {
	for _, e := range l.zones {
		e.Update()
	}
}

func (l *Level) GetStars(score *Score) {
	var i int
	for i = 0; i < 3; i++ {
		if float64(l.startimes[i]) < score.Time.Seconds() {
			break
		}
	}
	score.Stars = i
}
