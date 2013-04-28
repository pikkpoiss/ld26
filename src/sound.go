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
	m := C.mpg123_new(nil, nil)
	defer C.mpg123_delete(m)

	f := C.CString(file)

	if err := C.mpg123_open(m, f); err != C.MPG123_OK {
		panic("Error reading file")
	}
	defer C.mpg123_close(m)

	C.mpg123_scan(m)

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
