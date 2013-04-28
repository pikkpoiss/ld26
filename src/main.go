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
	"flag"
	"log"
	"runtime"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

func main() {
	var (
		err    error
		game   *Game
		window *twodee.Window
		system *twodee.System
	)
	var level = flag.Int("level", -1, "Skip to a level")
	flag.Parse()
	if system, err = twodee.Init(); err != nil {
		log.Fatalf("Couldn't init system: %v\n", err)
	}
	defer system.Terminate()
	window = &twodee.Window{Width: 1136, Height: 640, Scale: 1, Resize: false}
	if game, err = NewGame(system, window); err != nil {
		log.Fatalf("Couldn't start game: %v\n", err)
	}
	if *level != -1 {
		game.SetLevel(*level)
	}
	if err = game.Run(); err != nil {
		log.Fatalf("Exiting: %v\n", err)
	}
	log.Printf("Exiting peacefully")
}
