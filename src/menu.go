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
	index int
}

func NewMenuOption(s *twodee.Sprite, i int) *MenuOption {
	return &MenuOption{
		Sprite:     NewSprite(s),
		Selectable: false,
		index: i,
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

func (m *Menu) SetBounds(b twodee.Rectangle) {
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
