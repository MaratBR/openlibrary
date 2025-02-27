mkdir __vips
cd __vips

# install build system
sudo apt install -y \
    python3-mesonpy \
    build-essential pkg-config \
    libglib2.0-dev libexpat1-dev
# build deps
sudo apt install -y \
    libjpeg8 libjpeg8-dev \
    libpng16-16 libpng-dev \
    libwebp7 libwebp-dev


wget https://github.com/libvips/libvips/releases/download/v8.16.0/vips-8.16.0.tar.xz
tar xf vips-8.16.0.tar.xz
rm vips-8.16.0.tar.xz
cd vips-8.16.0
meson setup build
cd build
meson compile
meson test
sudo meson install