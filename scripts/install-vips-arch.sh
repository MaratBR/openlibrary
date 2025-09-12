mkdir __vips
cd __vips

# install build system
sudo pacman -Sy glibmm meson expat pkgconf glib2-devel libjpeg-turbo libpng libwebp libavif


wget https://github.com/libvips/libvips/releases/download/v8.16.0/vips-8.16.0.tar.xz
tar xf vips-8.16.0.tar.xz
rm vips-8.16.0.tar.xz
cd vips-8.16.0
meson setup build
cd build
meson compile
meson test
sudo meson install
