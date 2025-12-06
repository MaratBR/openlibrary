mkdir __vips
cd __vips

# Install build system dependencies
sudo dnf install -y \
    meson \
    gcc-c++ pkg-config \
    glib2-devel expat-devel

# Install build dependencies
sudo dnf install -y \
    libjpeg-turbo-devel \
    libpng-devel \
    libwebp-devel

# Download and extract vips
wget https://github.com/libvips/libvips/releases/download/v8.16.0/vips-8.16.0.tar.xz
tar xf vips-8.16.0.tar.xz
rm vips-8.16.0.tar.xz

cd vips-8.16.0

# Build using meson
meson setup build
cd build
meson compile
meson test
sudo meson install
