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

// #cgo CFLAGS: -I/Users/wesgoodman/Development/homebrew/include
// #cgo darwin LDFLAGS: -L/Users/wesgoodman/Development/homebrew/lib
// #cgo LDFLAGS: -lmpg123 -lao -ldl -lm -framework CoreFoundation
/*
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <errno.h>
#include "mpg123.h"
#include <ao/ao.h>
*/
import "C"

import (
	// "fmt"
	"unsafe"
)

type SoundSystem struct {
}

func NewSoundSystem() *SoundSystem {
	C.mpg123_init()
	C.ao_initialize()
	return &SoundSystem{}
}

func (s *SoundSystem) Shutdown() {
	C.ao_shutdown()
	C.mpg123_exit()
}

func (s *SoundSystem) DecodeTrack(file string, control chan int) {

	// var v1 *C.mpg123_id3v1
	// var v2 *C.mpg123_id3v2

	m := C.mpg123_new(nil, nil)
	defer C.mpg123_delete(m)

	f := C.CString(file)

	if err := C.mpg123_open(m, f); err != C.MPG123_OK {
		panic("Error reading file")
	}
	defer C.mpg123_close(m)

	C.mpg123_scan(m)
	// meta := C.mpg123_meta_check(m)

	// if meta == C.MPG123_ID3 && C.mpg123_id3(m, &v1, &v2) == C.MPG123_OK {
	// 	var title, artist, album, genre string
	// 	switch false {
	// 	case v2 == nil:
	// 		fmt.Println("ID3V2 tag found")
	// 		title = C.GoString(v2.title.p)
	// 		artist = C.GoString(v2.artist.p)
	// 		album = C.GoString(v2.album.p)
	// 		genre = C.GoString(v2.genre.p)

	// 	case v1 == nil:
	// 		fmt.Println("ID3V2 tag found")
	// 		title = C.GoString(&v1.title[0])
	// 		artist = C.GoString(&v1.artist[0])
	// 		album = C.GoString(&v1.album[0])
	// 		genre = "Unknown" // FIXME convert int to string
	// 	}

	// 	fmt.Println(title)
	// 	fmt.Println(artist)
	// 	fmt.Println(album)
	// 	fmt.Println(genre)
	// }

	default_driver := C.ao_default_driver_id()
	var format C.ao_sample_format
	var device *C.ao_device

	var channels, encoding C.int
	var rate C.long
	C.mpg123_getformat(m, &rate, &channels, &encoding)
	format.bits = 16
	format.channels = channels
	format.rate = C.int(rate)
	format.byte_format = C.AO_FMT_LITTLE

	device = C.ao_open_live(default_driver, &format, nil)
	if device == nil {
		panic("Error opening device")
		return
	}
	defer C.ao_close(device)

	var ret C.int
	var fill C.size_t
	buf := make([]C.uchar, 1024*16)

	for {
		ret = C.mpg123_read(m, (*C.uchar)(unsafe.Pointer(&buf)), 16*1024, &fill)
		if ret == C.MPG123_DONE {
			control <- 1
			break
		}
		C.ao_play(device, (*C.char)(unsafe.Pointer(&buf)), 16*1024)
	}
}
