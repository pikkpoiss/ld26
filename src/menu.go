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
	"fmt"
	"log"
	"sort"
)

const (
	MENU_DESELECTED = iota
	MENU_SELECTED
)

// Allow for sorting a list of sprites

type MenuOptions []*MenuOption

func (s MenuOptions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s MenuOptions) Len() int { return len(s) }

type MenuOptionsByX struct{ MenuOptions }

func (s MenuOptionsByX) Less(i, j int) bool {
	return s.MenuOptions[i].Bounds().Min.X < s.MenuOptions[j].Bounds().Min.X
}

type MenuOption struct {
	*Sprite
	Selectable bool
	index      int
}

func NewMenuOption(s *twodee.Sprite, i int) *MenuOption {
	return &MenuOption{
		Sprite:     NewSprite(s),
		Selectable: false,
		index:      i,
	}
}

func (o *MenuOption) SetState(state int) {
	o.Sprite.SetState(state)
	switch state {
	case MENU_DESELECTED:
		if o.Selectable {
			o.SetFrame(o.index)
		} else {
			o.SetFrame(o.index + 2)
		}
	case MENU_SELECTED:
		o.SetFrame(o.index + 1)
	}
}

// Interface to menus
type Navigatable interface {
	PrevSelection()
	NextSelection()
}

type Menu struct {
	system     *twodee.System
	background *Sprite
	options    []*MenuOption
	misc       []*Sprite
	index      int
}

func NewMenu(s *twodee.System) *Menu {
	return &Menu{
		system: s,
	}
}

func (m *Menu) Create(tileset string, index int, x, y, w, h float64) {
	switch tileset {
	case "background":
		s := m.system.NewSprite(tileset, x, y, w, h, index)
		m.background = NewSprite(s)
		m.background.SetTextureHeight(640)
		m.background.SetFrame(index)
	case "targets":
		s := NewMenuOption(m.system.NewSprite(tileset, x, y, w, h, index), index)
		s.SetState(MENU_DESELECTED)
		m.options = append(m.options, s)
	default:
		s := NewSprite(m.system.NewSprite(tileset, x, y, w, h, index))
		s.SetFrame(index)
		m.misc = append(m.misc, s)
	}
}

func (m *Menu) Draw() {
	if m.background != nil {
		m.background.Draw()
	}
	if m.misc != nil {
		for _, e := range m.misc {
			e.Draw()
		}
	}
	if m.options != nil {
		for _, e := range m.options {
			e.Draw()
		}
	}
}

func (m *Menu) Loaded(b twodee.Rectangle, props map[string]string) {
	sort.Sort(MenuOptionsByX{m.options})
	m.changeIndex(0)
	m.options[0].SetState(MENU_SELECTED)
}

func (m *Menu) changeIndex(i int) {
	var newindex = (i + len(m.options)) % len(m.options)
	if !m.options[newindex].Selectable {
		return
	}
	m.options[m.index].SetState(MENU_DESELECTED)
	m.index = newindex
	m.options[m.index].SetState(MENU_SELECTED)
}

func (m *Menu) SetSelectable(i int, selectable bool) {
	var index = (i + len(m.options)) % len(m.options)
	m.options[index].Selectable = selectable
	m.options[index].SetState(m.options[index].State())
}

func (m *Menu) NextSelection() {
	m.changeIndex(m.index + 1)
}

func (m *Menu) PrevSelection() {
	m.changeIndex(m.index - 1)
}

func (m *Menu) SetSelection(i int) {
	m.changeIndex(i)
}

func (m *Menu) GetSelection() int {
	return m.index
}

type LevelSelect struct {
	*Menu
	scores []*Score
	star   *Sprite
}

func NewLevelSelect(s *twodee.System, levels int) *LevelSelect {
	return &LevelSelect{
		Menu:   NewMenu(s),
		scores: make([]*Score, levels),
		star:   nil,
	}
}

func (l *LevelSelect) SetScores(scores []*Score) {
	l.scores = scores
}

func (l *LevelSelect) Loaded(b twodee.Rectangle, props map[string]string) {
	l.Menu.Loaded(b, props)
	l.star = NewSprite(l.Menu.system.NewSprite("select32", 0, 0, 2, 2, 0))
	l.star.SetFrame(1)
}

func (l *LevelSelect) Draw() {
	l.Menu.Draw()
	for i, s := range l.scores {
		if s == nil {
			continue
		}
		pt := l.Menu.options[i].Bounds().Min
		for j := 0; j < s.Stars; j++ {
			l.star.MoveTo(pt)
			pt.X += l.star.Width()
			l.star.Draw()
		}
	}
}

type Summary struct {
	system      *twodee.System
	background  []*Sprite
	stars       []*Sprite
	font        *twodee.Font
	pointTime   twodee.Point
	pointLevel  twodee.Point
	pointDamage twodee.Point
	score       *Score
	width       float64
	height      float64
	window      *twodee.Window
}

func NewSummary(s *twodee.System, font *twodee.Font, win *twodee.Window) *Summary {
	return &Summary{
		system: s,
		font:   font,
		window: win,
	}
}

func (s *Summary) Create(tileset string, index int, x, y, w, h float64) {
	switch tileset {
	case "options16":
		b := NewSprite(s.system.NewSprite(tileset, x, y, w, h, index))
		b.SetFrame(index)
		s.background = append(s.background, b)
	case "summary128":
		b := NewSprite(s.system.NewSprite(tileset, x, y, w, h, index))
		b.SetFrame(index)
		s.stars = append(s.stars, b)
	case "time":
		s.pointTime = twodee.Pt(x, y)
	case "level":
		s.pointLevel = twodee.Pt(x, y)
	case "damage":
		s.pointDamage = twodee.Pt(x, y)
	default:
		log.Printf("tileset: %v index: %v\n", tileset, index)
	}
}

func (s *Summary) textCoords(pt twodee.Point) twodee.Point {
	var (
		tx = pt.X / s.width * float64(s.window.Width)
		ty = pt.Y / s.height * float64(s.window.Height)
	)
	return twodee.Pt(tx, ty)
}

func (s *Summary) Draw() {
	if s.background != nil {
		for _, e := range s.background {
			e.Draw()
		}
	}
	if s.stars != nil {
		for _, e := range s.stars {
			e.Draw()
		}
	}
	var pt twodee.Point
	pt = s.textCoords(s.pointTime)
	s.font.Printf(pt.X, pt.Y, "%.1f seconds", s.score.Time.Seconds())
	pt = s.textCoords(s.pointLevel)
	s.font.Printf(pt.X, pt.Y, "Level %v", s.score.Level)
	pt = s.textCoords(s.pointDamage)
	str := fmt.Sprintf("Damage taken %.0f", s.score.Damage*100.0)
	s.font.Printf(pt.X, pt.Y, str)
}

func (s *Summary) SetMetrics(score *Score) {
	s.score = score
	for i := 0; i < len(s.stars); i++ {
		if i < s.score.Stars {
			s.stars[i].SetFrame(1)
		} else {
			s.stars[i].SetFrame(0)
		}
	}
}

func (s *Summary) Loaded(b twodee.Rectangle, props map[string]string) {
	s.width = b.Max.X - b.Min.X
	s.height = b.Max.Y - b.Min.Y
}
