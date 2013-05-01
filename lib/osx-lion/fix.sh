#!/usr/bin/env bash

function change {
  install_name_tool -change $1 $2 $3
}

function make_local {
  change $1/$2 @executable_path/$2 $3
}

PREFIX=/opt/twitter/lib
CPREFIX=/opt/twitter/Cellar/libvorbis/1.3.3/lib
GPREFIX=/usr/lib

chmod +w *.dylib
rm *.dylib

cp $GPREFIX/libGLEW.1.9.0.dylib .
cp $PREFIX/libSDL-1.2.0.dylib .
cp $PREFIX/libSDL_image-1.2.0.dylib .
cp $PREFIX/libSDL_mixer-1.2.0.dylib .
cp $PREFIX/libvorbis.dylib .
cp $PREFIX/libvorbis.0.dylib .
cp $PREFIX/libvorbisfile.dylib .
cp $PREFIX/libvorbisfile.3.dylib .

chmod +w *.dylib
make_local $GPREFIX libGLEW.1.9.0.dylib libGLEW.1.9.0.dylib
make_local $PREFIX libSDL-1.2.0.dylib libSDL-1.2.0.dylib
make_local $PREFIX libSDL-1.2.0.dylib libSDL_mixer-1.2.0.dylib
make_local $PREFIX libSDL-1.2.0.dylib libSDL_image-1.2.0.dylib
make_local $PREFIX libSDL_image-1.2.0.dylib libSDL_image-1.2.0.dylib
make_local $PREFIX libSDL_mixer-1.2.0.dylib libSDL_mixer-1.2.0.dylib
make_local $CPREFIX libvorbis.0.dylib libvorbisfile.dylib
make_local $CPREFIX libvorbis.0.dylib libvorbisfile.3.dylib
make_local $PREFIX libvorbisfile.3.dylib libvorbisfile.dylib
make_local $PREFIX libvorbisfile.3.dylib libvorbisfile.3.dylib
make_local $PREFIX libogg.0.dylib libvorbisfile.dylib
make_local $PREFIX libogg.0.dylib libvorbisfile.3.dylib
make_local $PREFIX libogg.0.dylib libvorbis.dylib
make_local $PREFIX libogg.0.dylib libvorbis.0.dylib
make_local $PREFIX libvorbis.0.dylib libvorbis.dylib
make_local $PREFIX libvorbis.0.dylib libvorbis.0.dylib
chmod -w *.dylib
otool -L *.dylib

