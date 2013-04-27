ld26
====

Ludum Dare 26

Theme: minimalism

Instructions
------------
Starting:

    git clone ...
    git submodule init
    git submodule update

Running:

    go run src/*.go

TODO
----

    [X] Wait for theme announcement
    [X] Init library and open window
    [ ] Bootleg initial sprite work
    [ ] Add load / swap level code
    [ ] Bootleg test level
    [ ] Handle button presses
    [ ] Add planetoid object code
    [ ] Collision detection
    [ ] Scoring system
    [ ] Lives system
    [ ] End game state / screen
    [ ] Intro screen
    [ ] Level select screen?  Level progress at least
    [ ] Production sprite work
    [ ] Determine "standard" level structures
    [ ] Build a track of levels of increasing difficulty

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
