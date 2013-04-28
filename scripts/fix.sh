install_name_tool -change /usr/lib/libGLEW.1.9.0.dylib @executable_path/libGLEW.1.9.0.dylib ld26
install_name_tool -change /usr/lib/libSDL_image-1.2.0.dylib @executable_path/libSDL_image-1.2.0.dylib ld26
otool -L ld26
