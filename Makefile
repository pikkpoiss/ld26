# Copyright 2012 Arne Roomann-Kurrik
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: build package clean run

PROJECT  = ld26
SOURCES  = $(wildcard src/*.go)

OSXLIBS  = $(wildcard lib/osx/*.dylib)
OSXBUILD = build/$(PROJECT)-osx/${PROJECT}.app/Contents

VERSION = $(shell cat VERSION)
REPLACE = s/9\.9\.9/$(VERSION)/g


clean:
	rm -rf build

$(OSXBUILD)/Info.plist: pkg/osx/Info.plist
	mkdir -p $(OSXBUILD)
	sed $(REPLACE) pkg/osx/Info.plist > $@

$(OSXBUILD)/MacOS/%.dylib: lib/osx/%.dylib
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/MacOS/launch.sh: scripts/launch.sh
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/MacOS/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go
	cd $(OSXBUILD)/MacOS/ && ../../../../../scripts/fix.sh

$(OSXBUILD)/Resources/%.icns: assets/%.icns
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/Resources/data/%: data/%
	mkdir -p $(dir $@)
	cp $< $@

build/$(PROJECT)-osx-$(VERSION).zip: \
	$(OSXBUILD)/Info.plist \
	$(subst lib/osx/,$(OSXBUILD)/MacOS/,$(wildcard lib/osx/*.dylib)) \
	$(OSXBUILD)/MacOS/launch.sh \
	$(OSXBUILD)/MacOS/$(PROJECT) \
	$(subst data/,$(OSXBUILD)/Resources/data/,$(wildcard data/*)) \
	$(subst assets/,$(OSXBUILD)/Resources/, $(wildcard assets/*.icns))
	cd build && zip -r $(notdir $@) $(PROJECT)-osx

build: build/$(PROJECT)-osx-$(VERSION).zip

run: build
	$(OSXBUILD)/MacOS/launch.sh
