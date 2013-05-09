ld26
====

Ludum Dare 26

Theme: minimalism

Ludum Dare entry:       http://www.ludumdare.com/compo/ludum-dare-26/?action=preview&uid=7913

Downloads, video, info: http://eg.regio.us/ld26/

Instructions (OSX)
------------------
Starting:

    brew install libogg libvorbis
    brew install sdl sdl_image
    brew install sdl_mixer --with-libvorbis  // Flag might be unnecessary.
    // Maybe `export LD_LIBRARY_PATH=/Users/wesgoodman/Development/homebrew/lib`
    git clone ...
    git submodule init
    git submodule update
    go get github.com/banthar/Go-SDL/mixer

Running:

    go run src/*.go

Instructions (Win)
------------------

Install:
	Install mingw64
	Install make
	http://sourceforge.net/projects/mingw-w64/files/External%20binary%20packages%20%28Win64%20hosted%29/make/
	http://sourceforge.net/apps/trac/mingw-w64/wiki/Make
	
	Place contents of bin in $PATH:
	http://ftp.gnome.org/pub/gnome/binaries/win32/gtk+/2.24/gtk+-bundle_2.24.10-20120208_win32.zip
	PKG_CONFIG_PATH = C:\MinGW64\lib\pkgconfig
	
	Download to a non-space path! C:\src
	http://www.libsdl.org/release/SDL-1.2.15.zip
	Make sure that make & dep DLLs is also in a non-space path!
	Make sure your PATH has no space paths in it!
	./configure
	
	http://www.libsdl.org/extras/win32/cross/README.txt
	http://www.libsdl.org/extras/win32/cross/cross-configure.sh
	http://www.libsdl.org/extras/win32/cross/cross-make.sh
	In non-Git Bash terminal:
	cross-configure.sh



--
	DLLs and headers
	http://www.libsdl.org/release/SDL-1.2.15-win32-x64.zip
	http://www.libsdl.org/release/SDL-devel-1.2.15-VC.zip
	http://www.libsdl.org/projects/SDL_mixer/release/SDL_mixer-1.2.12-win32-x64.zip
	http://www.libsdl.org/projects/SDL_mixer/release/SDL_mixer-devel-1.2.12-VC.zip
	
	https://groups.google.com/forum/#!topic/golang-nuts/CJofGE16KTk
	Create C:\MinGW64\lib\pkgconfig\sdl.pc
    # sdl pkg-config source file

    prefix=c:/MinGW64
    exec_prefix=${prefix}
    libdir=${exec_prefix}/lib
    includedir=${prefix}/include

    Name: sdl
    Description: Simple DirectMedia Layer is a cross-platform multimedia library designed to provide low level access to audio, keyboard, mouse, joystick, 3D hardware via OpenGL, and 2D video framebuffer.
    Version: 1.2.15
    Requires:
    Conflicts:
    Libs: -L${libdir}  -lmingw32 -lSDLmain -lSDL  -mwindows
    Libs.private: -lmingw32 -lSDLmain -lSDL  -mwindows  -liconv -lm -luser32 -lgdi32 -lwinmm -ldxguid
    Cflags: -I${includedir}/SDL -D_GNU_SOURCE=1 -Dmain=SDL_main

    go get -u github.com/banthar/Go-SDL/mixer
	go get -u github.com/banthar/Go-SDL/sdl

TODO
----

    [X] Wait for theme announcement
    [X] Init library and open window
    [X] Bootleg initial sprite work
    [X] Bootleg test level
    [X] Add planetoid object code
    [X] Game loop
    [X] Handle ESC/Close window
    [X] Add camera code
    [X] Load dev map and get running
    [X] Intro screen
    [X] Interface instead of relying on *Sprite for everything
    [X] Collision detection
    [X] Handle button presses
    [X] Player physics
    [X] Level select screen?  Level progress at least
    [X] Flags to skip intro / level select
    [X] Menus (exit, restart level, etc)
    [X] Don't require power of 2 image sizes
    [X] Make restart do something
    [X] Win a level
    [X] Add load / swap level code
    [X] End game state / screen
    [X] Scoring system
    [X] Show level scores on selection screen
    [X] Handle screen resize or prevent screen resize
    [X] Make blue regions hurt you instead of bouncing
    [X] Take damage if you're offscreen
    [X] Restart if the damage gets too high
    [X] Determine "standard" level structures
    [X] Build a track of levels of increasing difficulty
    [X] Code for music
    [X] Programmer music
    [X] Polished sprite work
    [X] Polished music
    [ ] Save / load level scores / progress
    [ ] Implement billboard system
    [ ] Add help billboards to menus and intro levels
    [ ] Code for sound effects
    [ ] Programmer sfx
    [ ] Polished sfx


Brainstorming
-------------
- Minimal artwork
- Simplified controls (one button)
- Emergence, making a world out of nothing
- "Minimalist gameshave small rulesets, narrow decision spaces, and
  abstractaudiovisual representations"
  http://www.academia.edu/1108618/Towards_Minimalist_Game_Design
- You're a planetoid
- Press one button to become attracted to the nearest object
- Get too close to the object you die
- Fly into outer space you die
- Have to keep pressing gravity button to orbit objects and change trajectory
- Puzzles like flying through a corridor of objects and needing to have gravity
  at the right time to avoid obstacles
- Make it obvious visually what you'd be attracted to
- Items which don't have mass / attractiveness
